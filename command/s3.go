package command

import (
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/reflection"
	"net/http"
	"time"
)

var (
	s3StandaloneOptions S3Options
)

type S3Options struct {
	filer                     *string
	bindIp                    *string
	port                      *int
	portGrpc                  *int
	config                    *string
	domainName                *string
	tlsPrivateKey             *string
	tlsCertificate            *string
	metricsHttpPort           *int
	allowEmptyFolder          *bool
	allowDeleteBucketNotEmpty *bool
	auditLogConfig            *string
	localFilerSocket          *string
	dataCenter                *string
}

func NewS3() *S3Options {
	return &S3Options{}
}

func (s3opt *S3Options) startS3Server() bool {
	filerBucketsPath := "/buckets"

	router := mux.NewRouter().SkipClean(true)
	var localFilerSocket string
	if s3opt.localFilerSocket != nil {
		localFilerSocket = *s3opt.localFilerSocket
	}
	s3ApiServer, s3ApiServer_err := s3api.NewS3ApiServer(router, &s3api.S3ApiServerOption{
		Filer:                     filerAddress,
		Port:                      *s3opt.port,
		Config:                    *s3opt.config,
		DomainName:                *s3opt.domainName,
		BucketsPath:               filerBucketsPath,
		GrpcDialOption:            grpcDialOption,
		AllowEmptyFolder:          *s3opt.allowEmptyFolder,
		AllowDeleteBucketNotEmpty: *s3opt.allowDeleteBucketNotEmpty,
		LocalFilerSocket:          localFilerSocket,
		DataCenter:                *s3opt.dataCenter,
	})
	if s3ApiServer_err != nil {
		glog.Fatalf("S3 API Server startup error: %v", s3ApiServer_err)
	}

	httpS := &http.Server{Handler: router}

	if *s3opt.portGrpc == 0 {
		*s3opt.portGrpc = 10000 + *s3opt.port
	}
	if *s3opt.bindIp == "" {
		*s3opt.bindIp = "localhost"
	}

	listenAddress := fmt.Sprintf("%s:%d", *s3opt.bindIp, *s3opt.port)
	s3ApiListener, s3ApiLocalListener, err := util.NewIpAndLocalListeners(*s3opt.bindIp, *s3opt.port, time.Duration(10)*time.Second)
	if err != nil {
		glog.Fatalf("S3 API Server listener on %s error: %v", listenAddress, err)
	}

	if len(*s3opt.auditLogConfig) > 0 {
		s3err.InitAuditLog(*s3opt.auditLogConfig)
		if s3err.Logger != nil {
			defer s3err.Logger.Close()
		}
	}

	// starting grpc server
	grpcPort := *s3opt.portGrpc
	grpcL, grpcLocalL, err := util.NewIpAndLocalListeners(*s3opt.bindIp, grpcPort, 0)
	if err != nil {
		glog.Fatalf("s3 failed to listen on grpc port %d: %v", grpcPort, err)
	}
	grpcS := pb.NewGrpcServer(security.LoadServerTLS(util.GetViper(), "grpc.s3"))
	s3_pb.RegisterSeaweedS3Server(grpcS, s3ApiServer)
	reflection.Register(grpcS)
	if grpcLocalL != nil {
		go grpcS.Serve(grpcLocalL)
	}

	if *s3opt.tlsPrivateKey != "" {
		glog.V(0).Infof("Start Seaweed S3 API Server %s at https port %d", util.Version(), *s3opt.port)
		if s3ApiLocalListener != nil {
			go func() {
				if err = httpS.ServeTLS(s3ApiLocalListener, *s3opt.tlsCertificate, *s3opt.tlsPrivateKey); err != nil {
					glog.Fatalf("S3 API Server Fail to serve: %v", err)
				}
			}()
		}
		if err = httpS.ServeTLS(s3ApiListener, *s3opt.tlsCertificate, *s3opt.tlsPrivateKey); err != nil {
			glog.Fatalf("S3 API Server Fail to serve: %v", err)
		}
	} else {
		glog.V(0).Infof("Start Seaweed S3 API Server %s at http port %d", util.Version(), *s3opt.port)
		if s3ApiLocalListener != nil {
			go func() {
				if err = httpS.Serve(s3ApiLocalListener); err != nil {
					glog.Fatalf("S3 API Server Fail to serve: %v", err)
				}
			}()
		}
		if err = httpS.Serve(s3ApiListener); err != nil {
			glog.Fatalf("S3 API Server Fail to serve: %v", err)
		}
	}

	return true

}

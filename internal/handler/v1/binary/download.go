package binary

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

type BinaryDownloader interface {
	BinaryDownload(ctx context.Context, data *storage.BinaryChunk) error
}

type BinaryDownloadFunc func(ctx context.Context, data *storage.BinaryChunk) error

func (f BinaryDownloadFunc) BinaryDownload(ctx context.Context, data *storage.BinaryChunk) error {
	return f(ctx, data)
}

var _ BinaryDownloader = BinaryDownloadFunc(nil)

type GRPCDownloadHandler func(req *pb.BidaryDownloadRequest, ds pb.GophKeeper_BinaryDownloadServer) error

func NewGRPCDownloadHandler(d BinaryDownloader) GRPCDownloadHandler {
	return func(req *pb.BidaryDownloadRequest, server pb.GophKeeper_BinaryDownloadServer) error {

		chSize := 4096
		var (
			done     bool
			err      error
			download pb.BinaryDownloadStream
		)
		data := storage.BinaryChunk{
			BinID: req.Id,
			Chunk: make([]byte, 0, chSize),
		}

		for !done {
			data.Offset += int64(len(data.Chunk))
			data.Chunk = data.Chunk[:0]
			err = d.BinaryDownload(server.Context(), &data)
			if err != nil || len(data.Chunk) == 0 {
				break
			} else if len(data.Chunk) < chSize {
				done = true
			}

			download.Chunk = data.Chunk
			err = server.Send(&download)
			if err != nil {
				break
			}
		}

		if errors.Is(err, storage.ErrNoContent) {
			return status.Error(codes.AlreadyExists, err.Error())
		} else if err != nil {
			logger.Errorf("download binary error: %w", err)
			return status.Error(codes.Internal, err.Error())
		}
		return nil
	}
}

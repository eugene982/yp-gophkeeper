package binary

import (
	"context"
	"io"

	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
)

type BinaryUploader interface {
	BinaryUpload(ctx context.Context, data storage.BinaryChunk) error
}

type BinaryUploadFunc func(ctx context.Context, data storage.BinaryChunk) error

func (f BinaryUploadFunc) BinaryUpload(ctx context.Context, data storage.BinaryChunk) error {
	return f(ctx, data)
}

var _ BinaryUploader = BinaryUploadFunc(nil)

type GRPCUploadHandler func(us pb.GophKeeper_BinaryUploadServer) error

func NewGRPCUploaderHandler(u BinaryUploader) GRPCUploadHandler {
	return func(server pb.GophKeeper_BinaryUploadServer) error {
		var (
			chunk  storage.BinaryChunk
			offset int64
		)

		var (
			err    error
			stream *pb.BinaryUplodStream
		)
		for err == nil {
			stream, err = server.Recv()
			if err == nil {
				chunk.BinID = stream.Id
				chunk.Chunk = stream.Chunk
				chunk.Offset = offset
				offset += int64(len(chunk.Chunk))
				err = u.BinaryUpload(server.Context(), chunk)
			} else {
				logger.Errorf("error upload binary: %w", err,
					"id", stream)
			}
		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			server.SendAndClose(&emptypb.Empty{})
			return err
		}
		return server.SendAndClose(&emptypb.Empty{})
	}
}

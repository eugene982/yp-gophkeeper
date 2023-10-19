package register

// func TestRPCRegisterHandler(t *testing.T) {

// 	tests := []struct {
// 		name       string
// 		login      string
// 		wantStatus codes.Code
// 	}{
// 		{name: "ok", login: "user", wantStatus: 0},
// 		{name: "already exists", login: "user", wantStatus: codes.AlreadyExists},
// 		{name: "internal error", login: "user", wantStatus: codes.Internal},
// 	}
// 	for _, tcase := range tests {
// 		t.Run(tcase.name, func(t *testing.T) {

// 			register := RegisterFunc(func(ctx context.Context, login, _ string) (string, error) {
// 				switch tcase.wantStatus {
// 				case 0:
// 					return login, nil
// 				case codes.AlreadyExists:
// 					return "", storage.ErrWriteConflict
// 				default:
// 					return "", fmt.Errorf("error")
// 				}
// 			})

// 			req := pb.RegisterRequest{
// 				Login:    tcase.login,
// 				Password: tcase.login,
// 			}

// 			resp, err := NewRPCRegisterHandler(register)(context.Background(), &req)
// 			if tcase.wantStatus == 0 {
// 				assert.NoError(t, err)
// 				require.NotNil(t, resp)
// 				assert.Equal(t, tcase.login, resp.Token)
// 			} else {
// 				assert.Error(t, err)
// 			}

// 		})
// 	}
// }

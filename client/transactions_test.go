package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/portto/aptos-go-sdk/models"
	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	impl := NewAptosClient("https://fullnode.mainnet.aptoslabs.com/v1")
	p := models.ViewParam{
		Function: "0x0e69e1d1069f086aca14daccbd3183848a1a446f5c3d3ea09bfa964e9324798c::BetterSwaps2::fetch_res",
		Arguments: []string{"1"},
		TypeArguments: []string{},
	}
	_, err := impl.GetAccountResources(context.Background(), "0x0e69e1d1069f086aca14daccbd3183848a1a446f5c3d3ea09bfa964e9324798c")
	t.Log(err)
	resp := []interface{}{}
	err = impl.View(context.Background(), p, &resp, nil)
	assert.NoError(t, err)
	t.Logf("%v", resp)
}

func TestGetTransactionByHash(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		resp := []byte(mockTx)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			n, err := w.Write(resp)
			assert.NoError(t, err)
			assert.Equal(t, len(resp), n)
		}))
		impl := NewAptosClient(srv.URL)
		tx, err := impl.GetTransactionByHash(mockCTX, mockTxHash)
		assert.NoError(t, err)
		assert.Equal(t, "0x"+mockTxHash, tx.Hash)
	})
}

func TestWaitForTransaction(t *testing.T) {
	t.Run("Retry", func(t *testing.T) {
		resp := []byte(mockTx)
		var count int
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if count == 0 {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, err := w.Write([]byte("bad request"))
				assert.NoError(t, err)
				count += 1
			} else {
				n, err := w.Write(resp)
				assert.NoError(t, err)
				assert.Equal(t, len(resp), n)
			}
		}))
		impl := NewAptosClient(srv.URL)
		err := impl.WaitForTransaction(mockCTX, mockTxHash)
		assert.NoError(t, err)
	})
}

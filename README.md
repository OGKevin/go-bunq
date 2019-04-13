bunq package is a bunq client for the bunq api: https://doc.bunq.com

Be aware that there could be unexported methods, structs etc that should be exported
as this package was literally extracted from another package, and some custom code has been removed.

Just ask if you are in doubt.

Maybe in the future ill add a generator, if not, all endpoints need to be added by hand. However that is not that hard anyway. 

```go
package main

import "github.com/OGKevin/go-bunq/bunq"
import "context"
import "log"

func main() {
    	
  key, err := bunq.CreateNewKeyPair()
  if err != nil {
      panic(err)
  }

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
  c := bunq.NewClient(ctx, bunq.BaseURLSandbox, key, "sandbox_ab7df7985a66133b1abecf42871801edaafe5bc51ef9769f5a032876")
  err = c.Init()
  if err != nil {
      panic(err)
  }
  
  for i := 0; i < 20; i++ {
      log.Print(i)
  
      _, err = c.PaymentService.CreatePaymentBatch(
          10111,
          bunq.PaymentBatchCreate{
              Payments: generateBatchEntries(100),
          },
      )
  if err != nil {
    panic(err)
    }
  }
}

func generateBatchEntries(nr int) []bunq.PaymentCreate {
	var entries []bunq.PaymentCreate

	for i := 0; i < nr; i++ {
		entries = append(
			entries,
			bunq.PaymentCreate{
				Amount: bunq.Amount{
					Currency: "EUR",
					Value:    "0.01",
				},
				CounterpartyAlias: bunq.Pointer{
					PType: "EMAIL",
					Value: "bravo@bunq.com",
				},
				Description: "test",
			},
		)
	}

	return entries
}

```
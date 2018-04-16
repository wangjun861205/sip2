package sip2

// func TestRequest(t *testing.T) {
// 	pool, err := NewClientPool("112.230.195.36", 2030, 20, 10, 5, true)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	patronStatus, err := pool.QueryPatronStatus(Chinese, "", "", "", "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	b, err := json.Marshal(patronStatus)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("%#v\n", string(b))
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	err = pool.Close(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

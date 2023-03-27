package gormer

// func TestFactory(t *testing.T) {
// 	db := New(&Config{
// 		Default: "db1",
// 		Connections: map[string]ConnectionFunc{
// 			"db1": func() *gorm.DB {
// 				return nil
// 			},
// 			"db2": func() *gorm.DB {
// 				return nil
// 			},
// 		},
// 	})
//
// 	SetInstance(db)
//
// 	Factory().Connection("db1").First(nil)
// }

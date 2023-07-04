package bots

import (
	"github.com/donglei1234/platform/services/common/nosql/document"
	errors2 "github.com/donglei1234/platform/services/common/nosql/errors"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/bot"
	"github.com/donglei1234/platform/services/common/dupblock"
	"github.com/donglei1234/platform/services/common/nosql/document/badger"
	"github.com/donglei1234/platform/services/common/utils"
)

func NosqlSmokeTest(config *bot.Config) (err error) {
	ErrUnexpectedResultLength := errors.New("ErrUnexpectedResultLength")

	time.Sleep(config.StartDelay)
	startTime := time.Now()
	if dir, err := utils.NewTempDir("nosql_bot_test_"); err != nil {
		if bot.LogErrorAndCheckToBreak(err, config, "NewTempDir") {
			return err
		}
	} else {
		defer func() {
			if err := dir.Cleanup(); err != nil {
				config.Logger.Error("dir.Cleanup", zap.Error(err))
			}
		}()

		for {
			config.Logger.Info("		       ┌───────────────────────────────┐")
			config.Logger.Info("		       │        NoSQL Test Suite       │")
			config.Logger.Info("		       └───────────────────────────────┘")

			// create the document store provider in question
			var provider document.DocumentStoreProvider
			config.Logger.Info("Creating a new Document Store Provider ...")
			if provider, err = badger.NewDocumentStoreProvider(dir.Path(), 5*time.Minute, config.Logger); err != nil {
				if bot.LogErrorAndCheckToBreak(err, config, "NewDocumentStoreProvider") {
					break
				}
			}
			// open the provider and run the test suite against it
			config.Logger.Info("Opening the Document Store ...")
			var store document.DocumentStore
			if store, err = provider.OpenDocumentStore("smokeTest"); err != nil {
				if bot.LogErrorAndCheckToBreak(err, config, "OpenDocumentStore") {
					break
				}
			} else {
				config.Logger.Info("provider.OpenDocumentStore", zap.String("Store:", store.Name()))
				var err error

				config.Logger.Info("Creating a new Key '/test/keys' ...")
				var testKey document.Key
				if testKey, err = document.NewKeyFromString("/test/keys"); err != nil {
					if bot.LogErrorAndCheckToBreak(err, config, "NewKeyFromString") {
						break
					}
				}
				config.Logger.Info("nosql.NewKeyFromString", zap.String("Key:", testKey.String()))

				var testVersion document.Version
				src := map[string]interface{}{"a": "b"}

				config.Logger.Info("Setting the value of the keys to {`a`:`b`} ...",
					zap.String("Key:", testKey.String()))

				if testVersion, err = store.Set(
					testKey,
					document.WithAnyVersion(),
					document.WithTTL(1*time.Minute),
					document.WithSource(src)); err != nil {
					if bot.LogErrorAndCheckToBreak(err, config, "store.Set") {
						break
					}
				} else {
					config.Logger.Info("store.Set", zap.Any("Version:", testVersion))
					config.Logger.Info("Checking if the Document Store now contains the Key ...")
					if ok, err := store.Contains(testKey); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Contains") {
							break
						}
					} else if !ok {
						if bot.LogErrorAndCheckToBreak(
							errors2.ErrKeyNotFound,
							config,
							"store.Contains",
						) {
							break
						}
					} else {
						config.Logger.Info("store.Contains", zap.Bool("Found:", ok))
					}

					config.Logger.Info("Getting the keys's value to confirm it was set ...")
					var dst interface{}
					if testVersion, err = store.Get(
						testKey,
						document.WithVersion(testVersion),
						document.WithDestination(&dst)); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Get") {
							break
						}
					}
					config.Logger.Info(
						"store.Get",
						zap.String("Key:", testKey.String()),
						zap.Any("Version:", testVersion),
					)

					config.Logger.Info("Applying DUPBlock test cases to the document store ...")
					testCases := []string{
						"keys /test/keys\nset test\n5\n*",
						"keys /test/keys\nset test2\n0\n*",
						"keys /test/keys\ncpy test test2\n",
						"keys /test/keys\nset test2\n10\n*",
						"keys /test/keys\nmov test2 test\n",
					}
					for _, tc := range testCases {
						if reader, err := dupblock.NewTextReader(dupblock.WithBytes([]byte(tc))); err != nil {
							if bot.LogErrorAndCheckToBreak(
								err,
								config,
								"dupblock.NewTextReader",
							) {
								break
							}
						} else if err := store.ApplyDUPBlock(reader); err != nil {
							if bot.LogErrorAndCheckToBreak(
								err,
								config,
								"store.ApplyDUPBlock",
							) {
								break
							}
						}
					}
					config.Logger.Info("All DUPBlock test cases were applied!")

					// prep data for ListKeys...
					config.Logger.Info("Inserting dataset for ListKeys & Scan tests ...")
					const (
						expectedTestKeys = 8 // Keys that begin with `/testkeys/`
					)
					for k, v := range map[string]interface{}{
						"/test/list_test_1":      map[string]interface{}{"Idx": "a", "Index": "1"},
						"/test/list_test_2":      map[string]interface{}{"Idx": "b", "value": "2"},
						"/testkeys/list_test_3":  map[string]interface{}{"Idx": "c", "int": 3},
						"/testkeys/list_test_4":  map[string]interface{}{"Idx": "d", "flag": false},
						"/testkeys/scan/test_5":  map[string]interface{}{"Idx": "e", "User": "Guybrush Threepwood"},
						"/testkeys/scan/test_6":  map[string]interface{}{"Idx": "f", "User": "Murray"},
						"/testkeys/scan/test_7":  map[string]interface{}{"Idx": "g", "User": "Murray"},
						"/testkeys/scan/test_8":  map[string]interface{}{"Idx": "h", "Value": 1.0},
						"/testkeys/scan/test_9":  map[string]interface{}{"Idx": "i", "Value": 2},
						"/testkeys/scan/test_10": map[string]interface{}{"Idx": "j", "Value": 3.14},
					} {

						config.Logger.Info("Setting a value ...", zap.String("Key:", k), zap.Any("Value:", v))

						var version document.Version
						if version, err = store.Set(
							document.NewKeyFromStringUnchecked(k),
							document.WithAnyVersion(),
							document.WithTTL(1*time.Minute),
							document.WithSource(v),
						); err != nil {
							if bot.LogErrorAndCheckToBreak(err, config, "store.Set") {
								break
							}
						} else {
							config.Logger.Info("store.Set", zap.Any("Version", version))
						}
					}

					// and quickly verify that the last keys from the set was added successfully
					if ok, err := store.Contains(
						document.NewKeyFromStringUnchecked("/testkeys/list_test_4"),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Contains") {
							break
						}
					} else if !ok {
						if bot.LogErrorAndCheckToBreak(
							errors2.ErrKeyNotFound,
							config,
							"store.Contains",
						) {
							break
						}
					} else {
						config.Logger.Info("store.Contains", zap.Bool("Found:", ok))
					}

					// test ListKeys
					config.Logger.Info("Listing keys that start with `/testkeys/` ...")
					if keyList, err := store.ListKeys("/testkeys/", document.WithNoLimit()); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.ListKeys") {
							break
						}
					} else if len(keyList) != expectedTestKeys {
						config.Logger.Error("ListKeys returned an unexpected result length.",
							zap.Int("Length:", len(keyList)), zap.Int("Expected:", expectedTestKeys))
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.ListKeys",
						) {
							break
						}
					} else {
						for _, key := range keyList {
							config.Logger.Info("", zap.String("Result:", key.String()))
						}
						config.Logger.Info("Successfully listed keys matching prefix `/testkeys/`!")
					}

					config.Logger.Info("Listing keys that start with" +
						" `/look-out-behind-you-it-s-a-three-headed-monkey/` ...")
					if keyList, err := store.ListKeys(
						"/look-out-behind-you-it-s-a-three-headed-monkey/",
						document.WithNoLimit(),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.ListKeys") {
							break
						}
					} else if len(keyList) != 0 {
						config.Logger.Error("ListKeys returned an unexpected result length.",
							zap.Int("Length:", len(keyList)), zap.Int("Expected:", 0))
						for _, key := range keyList {
							config.Logger.Info("", zap.String("Result:", key.String()))
						}
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.ListKeys",
						) {
							break
						}
					} else {
						config.Logger.Info("Successfully listed 0 keys matching nonexistent prefix.")
					}

					// validate ListKeys with limit and offset
					if keyList, err := store.ListKeys("/testkeys/", document.WithLimit(3)); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.ListKeys") {
							break
						}
					} else if len(keyList) != 3 {
						config.Logger.Error("ListKeys returned an unexpected result length.",
							zap.Int("Length:", len(keyList)), zap.Int("Expected:", 3))
						for _, key := range keyList {
							config.Logger.Info("", zap.String("Result:", key.String()))
						}
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.ListKeys",
						) {
							break
						}
					} else {
						config.Logger.Info("Listing keys that start with `/testkeys/`, limit 3 ...")
						for _, key := range keyList {
							config.Logger.Info("", zap.String("Result:", key.String()))
						}
						config.Logger.Info("Successfully listed keys matching prefix `/testkeys/`!")
					}

					config.Logger.Info("Listing keys that start with `/testkeys/`, offset 3 limit 2 ...")
					if keyList, err := store.ListKeys(
						"/testkeys/",
						document.WithOffset(3),
						document.WithLimit(2),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.ListKeys") {
							break
						}
					} else if len(keyList) != 2 {
						config.Logger.Error("ListKeys returned an unexpected result length.",
							zap.Int("Length:", len(keyList)), zap.Int("Expected:", 2))
						for _, key := range keyList {
							config.Logger.Info("", zap.String("Result:", key.String()))
						}
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.ListKeys",
						) {
							break
						}
					} else {
						config.Logger.Info("Listing keys that start with `/testkeys/`, offset 3 limit 2 ...")
						for _, key := range keyList {
							config.Logger.Info("", zap.String("Result:", key.String()))
						}
						config.Logger.Info("Successfully listed keys matching prefix `/testkeys/`!")
					}

					// test Scan
					type scanTest struct {
						Idx   string
						User  string
						Value interface{}
					}
					scanDest := scanTest{}
					config.Logger.Info("Scanning for single existant user ...")
					reqUser := "Guybrush Threepwood"
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanDest),
						document.MatchKeyValue("User", reqUser),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != 1 {
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else if scanDest.User != reqUser {
						if bot.LogErrorAndCheckToBreak(
							errors.New("ErrResultMismatch"),
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info("Found record:", zap.String("User:", reqUser))
					}

					config.Logger.Info("Scanning for nonexistant user ...")
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanDest),
						document.MatchKeyValue("User", "LeChuck"),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != 0 {
						// we're not expecting a response, freak out
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					}

					config.Logger.Info("Scanning for multi-entry user ...")
					reqUser = "Murray"
					expectedRows := 2
					scanRes := make([]scanTest, 2, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchKeyValue("User", reqUser),
						document.WithLimit(expectedRows),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.String("User", reqUser),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							config.Logger.Info(" ", zap.String("Result Row:", row.Idx))
						}
					}

					// test WithOffset
					config.Logger.Info("Scanning with offset...")
					expectedRows = 1
					scanRes = make([]scanTest, 1, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchKeyValue("User", reqUser),
						document.WithOffset(1),
						document.WithLimit(expectedRows),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there is one entry for this user, after offset
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.String("User", reqUser),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							config.Logger.Info(" ", zap.String("Result Row:", row.Idx))
						}
					}

					// test WithKeyLike
					config.Logger.Info("Scanning for users matching pattern...")
					expectedRows = 2
					reqUser = "Mu%"
					scanRes = make([]scanTest, 2, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchKeyLike("User", reqUser),
						document.WithNoLimit(),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.String("User", reqUser),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							config.Logger.Info(" ", zap.String("Result Row:", row.Idx))
						}
					}

					// test unset scan type
					config.Logger.Info("Scanning with unset type...")
					expectedRows = 0
					scanRes = make([]scanTest, 0, 5)
					if num, err := store.Scan("/testkeys/scan/", document.WithDestination(&scanRes)); err == nil {
						if bot.LogErrorAndCheckToBreak(
							errors.Wrap(
								errors2.ErrInternal, "got nil instead of expected error"),
							config,
							"store.Scan",
						) {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else if err != errors2.ErrNoScanType && errors.Cause(err) != errors2.ErrNoScanType {
						if bot.LogErrorAndCheckToBreak(
							errors.Wrap(errors2.ErrInternal, err.Error()),
							config,
							"store.Scan",
						) {
							break
						}
					}

					// test multiple chained queries
					config.Logger.Info("Scanning for multiple conditions ...")
					reqUser = "Murray"
					expectedRows = 1
					scanRes = make([]scanTest, 1, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchKeyValue("User", reqUser),
						document.MatchKeyValue("Idx", "g"),
						document.WithNoLimit(),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.String("User", reqUser),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							if row.Idx != "" {
								config.Logger.Info("", zap.Any("Row:", row))
							}
						}
					}

					// test multiple chained queries
					config.Logger.Info("Scanning for numerical equality ...")
					reqNum := 2.0
					expectedRows = 1
					scanRes = make([]scanTest, 1, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchNumber("Value", document.ScanOpEquals, reqNum),
						document.WithNoLimit(),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.Float64("reqNum", reqNum),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							if row.Idx != "" {
								config.Logger.Info(
									"",
									zap.String("Row:", row.Idx),
									zap.Float64("Row Value:", row.Value.(float64)),
								)
							}
						}
					}

					config.Logger.Info("Scanning for numerical < ...")
					reqNum = 3
					expectedRows = 2
					scanRes = make([]scanTest, 2, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchNumber("Value", document.ScanOpLessThan, reqNum),
						document.WithNoLimit(),
					); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.Float64("reqNum", reqNum),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							if row.Idx != "" {
								config.Logger.Info(
									"",
									zap.String("Row:", row.Idx),
									zap.Float64("Row Value:", row.Value.(float64)),
								)
							}
						}
					}

					config.Logger.Info("Scanning for numerical > ...")
					reqNum = 3.1
					expectedRows = 1
					scanRes = make([]scanTest, 1, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchNumber("Value", document.ScanOpGreaterThan, reqNum),
						document.WithNoLimit()); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.Float64("reqNum", reqNum),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							if row.Idx != "" {
								config.Logger.Info(
									"",
									zap.String("Row:", row.Idx),
									zap.Float64("Row Value:", row.Value.(float64)),
								)
							}
						}
					}

					// test MatchRegex
					regex := "e+pw.o"
					config.Logger.Info("Scanning for regex...", zap.String("Regex", regex))
					expectedRows = 1
					scanRes = make([]scanTest, 1, 5)
					if num, err := store.Scan(
						"/testkeys/scan/",
						document.WithDestination(&scanRes),
						document.MatchRegex(regex),
					); errors.Cause(err) == errors2.ErrDriverFailure {
						config.Logger.Info("Received anticipated error", zap.String("Error", err.Error()))
					} else if err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "store.Scan") {
							break
						}
					} else if num != expectedRows {
						// there are two entries for this user
						if bot.LogErrorAndCheckToBreak(
							ErrUnexpectedResultLength,
							config,
							"store.Scan",
						) {
							break
						}
					} else {
						config.Logger.Info(
							"Scan response",
							zap.Int("Expected", expectedRows),
							zap.String("regex", regex),
							zap.Int("Results", num),
						)
						for _, row := range scanRes {
							config.Logger.Info(
								"",
								zap.String("Row:", row.Idx),
								zap.String("Row User:", row.User),
							)
						}
					}

					// test Remove
					for k, v := range map[string]interface{}{
						"/test/remove_test_any": map[string]interface{}{"Index": "delete_me"},
					} {
						config.Logger.Info("Setting a value", zap.String("Key", k), zap.Any("Value", v))

						key := document.NewKeyFromStringUnchecked(k)
						if _, err = store.Set(
							key,
							document.WithAnyVersion(),
							document.WithTTL(10*time.Minute),
							document.WithSource(v),
						); err != nil {
							if bot.LogErrorAndCheckToBreak(err, config, "store.Set") {
								break
							}
						}

						// make sure it was set
						if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyNotFound,
								config,
								"store.Contains",
							) {
								break
							}
						} else if !ok {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyNotFound,
								config,
								"store.Contains",
							) {
								break
							}
						}

						// remove it
						config.Logger.Info("Removing keys", zap.String("Key ID", k))
						if err := store.Remove(key, document.WithAnyVersion()); err != nil {
							if bot.LogErrorAndCheckToBreak(err, config, "store.Remove") {
								break
							}
						}

						// make sure it has been removed
						if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyExists,
								config,
								"store.Contains",
							) {
								break
							}
						} else if ok {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyExists,
								config,
								"store.Contains",
							) {
								break
							}
						}
						break
					}

					// test Remove (with a version provided) - Yes, this is mostly redundant and could be consolidated
					for k, v := range map[string]interface{}{
						"/test/remove_test_ver": map[string]interface{}{"Index": "delete_me_too"},
					} {
						config.Logger.Info("Setting a value", zap.String("Key", k), zap.Any("Value", v))

						key := document.NewKeyFromStringUnchecked(k)
						var version document.Version
						if version, err = store.Set(
							key,
							document.WithAnyVersion(),
							document.WithTTL(10*time.Minute),
							document.WithSource(v),
						); err != nil {
							if bot.LogErrorAndCheckToBreak(err, config, "store.Set") {
								break
							}
						}

						// make sure it was set
						if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyNotFound,
								config,
								"store.Contains",
							) {
								break
							}
						} else if !ok {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyNotFound,
								config,
								"store.Contains",
							) {
								break
							}
						}

						var dst interface{}
						var getVersion document.Version
						if getVersion, err = store.Get(
							key,
							document.WithAnyVersion(),
							document.WithDestination(&dst),
						); err != nil {
							if bot.LogErrorAndCheckToBreak(err, config, "store.Get") {
								break
							}
						} else if version < getVersion {
							/* In some cases (Badger) it is expected for Get() to return a higher version than the immediately
							 * preceding Set() call (because Badger internals). So we only log an warning about this undesirable
							 * case at the time being.
							 */
							config.Logger.Warn(
								"Warning: Set returned a different version than Get.",
								zap.Uint64("SetVersion", version),
								zap.Uint64("GetVersion", getVersion),
							)
							version = getVersion
						} else if version != getVersion {
							/* However, if the version number ever rolls backward? That's very bad.
							 */
							if bot.LogErrorAndCheckToBreak(
								errors.Wrap(errors2.ErrInvalidVersioning, "unexpected version rollback"),
								config,
								"store.Get",
							) {
								break
							}
						} else {
							config.Logger.Info("store.Get", zap.Any("Version:", version))
						}

						// remove it
						config.Logger.Info("Removing the keys ...",
							zap.String("Key:", k), zap.Any("Version", version))
						if err := store.Remove(key, document.WithVersion(version)); err != nil {
							if bot.LogErrorAndCheckToBreak(err, config, "store.Remove") {
								break
							}
						}

						// make sure it has been removed
						if ok, err := store.Contains(document.NewKeyFromStringUnchecked(k)); err != nil {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyExists,
								config,
								"store.Contains",
							) {
								break
							}
						} else if ok {
							if bot.LogErrorAndCheckToBreak(
								errors2.ErrKeyExists,
								config,
								"store.Contains",
							) {
								break
							}
						} else {
							config.Logger.Info("store.Contains", zap.Bool("Found:", ok))
							config.Logger.Info("Successfully removed the keys!")
						}
						break
					}
					// clean up before displaying a "pass"
					if err = provider.Shutdown(); err != nil {
						if bot.LogErrorAndCheckToBreak(err, config, "provider.Shutdown") {
							break
						}
					}
					// make sure the cleanup doesn't panic
					provider = nil

					config.Logger.Info("		       ┌───────────────────────────────┐")
					config.Logger.Info("		       │    NoSQL Test Result: PASS    │")
					config.Logger.Info("		       └───────────────────────────────┘")
				}
			}

			if !config.Repeat || time.Now().Sub(startTime) > config.Duration {
				break
			} else {
				time.Sleep(config.RepeatDelay)
			}
		}
	}
	return nil
}

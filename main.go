// ____                       ____
// / ___| _   _ _ __   ___    / ___| _ __   __ _  ___ ___
// \___ \| | | | '_ \ / __|___\___ \| '_ \ / _` |/ __/ _ \
//  ___) | |_| | | | | (_|_____|__) | |_) | (_| | (_|  __/
// |____/ \__, |_| |_|\___|   |____/| .__/ \__,_|\___\___|
//        |___/                     |_|
//
//                      .';codxkOOkxdoc;'.
//                   ,lxKNMMMMWWWWWWMMMWN0xl'
//                .lONMMWKkdl:;,''',:cokKWMMNOl,..
//              .lKWMWKo;.              .,lONMMWXK0Okdl;'.
//             ,OWMWO:.                     ,xNMMMWNXWMMN0d;.
//            :XMMKc.                         ;0WMWO;,cd0WMNO,
//           :XMWO'                            .kWMWd.  .oNMM0,
//        .:kXMMO'                              .kWMNc   '0MMX:
//      ,dKWMMMX;                                ;XMMO. .oNMWk.
//    ,kNMNKXMMk.                                .kMMNc'xNMWk.
//  .oNMWO;.dMMo                                  dMMWXXWW0:.
// .dWMXl.  lWMd                                 ,0MMMMNk:.
// ;XMMx.   :NM0'                           .':dONMMMNx,
// ;XMM0,   .xMWo.          .oxoc:,..     .;kNMMWWMMWx.
// .lNMMXx:'.;KMNo.          .OMMMWX0ko,   .x0xloKMMK,
//   ,dKWMMN0OXMMNOdooooddxkONMMMMMMMWO,    .. ,0WMK;
//     .;cdk0XXWMMMMMMMMMWWNXXWMMMMW0c.      .cKMW0,
//          ...,xNMMMXd:;,     dWN0c.      .cOWMXo.
//               ,dXMWXxc'.   'l:'.    .':dKWWKd'
//                 .lkXWMNKkdolcc:clodk0NMWKx:.
//                    .;oxOKNWWWWWWWNX0kdc,.
//                         ..,,,,,,,'..
//
// Authors:
// Nathan Laney, Dylan Halstead
// 2023

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
	"github.com/Sync-Space-49/syncspace-server/routers"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Err(err).Msg("failed to run server")
	}
}

func run() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	db, err := db.New(cfg.DB.DBUser, cfg.DB.DBPass, cfg.DB.DBURI, cfg.DB.DBName)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    cfg.APIHost,
		Handler: routers.NewAPI(cfg, db),
	}

	fmt.Printf("Running the server on %s\n", cfg.APIHost)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed while running server: %w", err)
	}

	return nil
}

// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwar
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exportimporters

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/exposure-notifications-server/internal/admin"
	exportimportdatabase "github.com/google/exposure-notifications-server/internal/exportimport/database"
	exportimportmodel "github.com/google/exposure-notifications-server/internal/exportimport/model"
	"github.com/google/exposure-notifications-server/internal/serverenv"
)

type viewController struct {
	config *admin.Config
	env    *serverenv.ServerEnv
}

func NewView(c *admin.Config, env *serverenv.ServerEnv) admin.Controller {
	return &viewController{config: c, env: env}
}

func (v *viewController) Execute(c *gin.Context) {
	ctx := c.Request.Context()

	db := exportimportdatabase.New(v.env.Database())
	model := new(exportimportmodel.ExportImport)

	if idRaw := c.Param("id"); idRaw != "" && idRaw != "0" {
		id, err := strconv.ParseInt(idRaw, 10, 64)
		if err != nil {
			admin.ErrorPage(c, fmt.Sprintf("Failed to parse `id` param: %s", err))
			return
		}

		model, err = db.GetConfig(ctx, id)
		if err != nil {
			admin.ErrorPage(c, fmt.Sprintf("Failed to load export importer config: %s", err))
			return
		}
	}

	m := make(admin.TemplateMap)
	m["model"] = model
	c.HTML(http.StatusOK, "export-importer", m)
	c.Abort()
}

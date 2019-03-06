/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package synchronizer

import (
	"context"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/auth_synchronizer/app/options"
	"configcenter/src/scene_server/auth_synchronizer/pkg/synchronizer/meta"
	"configcenter/src/storage/dal"
)

// AuthSynchronizer stores all related resource
type AuthSynchronizer struct {
	Config *options.Config
	*backbone.Engine
	db          dal.RDB
	ctx         context.Context
	Workers     *[]Worker
	WorkerQueue chan meta.WorkRequest
	Producer    *Producer
}

// NewSynchronizer new a synchronizer object
func NewSynchronizer(ctx context.Context, config *options.Config, backbone *backbone.Engine) *AuthSynchronizer {
	return &AuthSynchronizer{ctx: ctx, Config: config, Engine: backbone}
}

// Run do start synchronizer
func (d *AuthSynchronizer) Run() error {
	blog.Infof("auth synchronizer start...")

	// init queue
	d.WorkerQueue = make(chan meta.WorkRequest, 1000)

	// make fake handler
	handler := new(meta.FakeHander)

	// init worker
	workers := make([]Worker, 3)
	for w := 1; w <= 3; w++ {
		worker := NewWorker(w, d.WorkerQueue, handler)
		workers = append(workers, *worker)
		worker.Start()
	}
	d.Workers = &workers

	// init producer
	d.Producer = NewProducer(d.WorkerQueue)
	blog.Infof("auth synchronizer started")
	return nil
}

/*
 *
 * Copyright 2019 The KubeSphere Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package openpitrix

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"kubesphere.io/kubesphere/pkg/models"
	"kubesphere.io/kubesphere/pkg/server/params"
	cs "kubesphere.io/kubesphere/pkg/simple/client"
	"kubesphere.io/kubesphere/pkg/simple/client/openpitrix"
	"openpitrix.io/openpitrix/pkg/pb"
)

func CreateCategory(request *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	op, err := cs.ClientSets().OpenPitrix()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	r := &pb.CreateCategoryRequest{
		Name:        &wrappers.StringValue{Value: request.Name},
		Locale:      &wrappers.StringValue{Value: request.Locale},
		Description: &wrappers.StringValue{Value: request.Description},
	}
	if request.Icon != nil {
		r.Icon = &wrappers.BytesValue{Value: request.Icon}
	}

	resp, err := op.Category().CreateCategory(openpitrix.SystemContext(), r)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return &CreateCategoryResponse{
		CategoryId: resp.GetCategoryId().GetValue(),
	}, nil
}

func DeleteCategory(id string) error {
	op, err := cs.ClientSets().OpenPitrix()
	if err != nil {
		klog.Error(err)
		return err
	}
	_, err = op.Category().DeleteCategories(openpitrix.SystemContext(), &pb.DeleteCategoriesRequest{
		CategoryId: []string{id},
	})
	if err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

func PatchCategory(id string, request *ModifyCategoryRequest) error {
	op, err := cs.ClientSets().OpenPitrix()
	if err != nil {
		klog.Error(err)
		return err
	}
	modifyCategoryRequest := &pb.ModifyCategoryRequest{
		CategoryId: &wrappers.StringValue{Value: id},
	}
	if request.Name != nil {
		modifyCategoryRequest.Name = &wrappers.StringValue{Value: *request.Name}
	}
	if request.Locale != nil {
		modifyCategoryRequest.Locale = &wrappers.StringValue{Value: *request.Locale}
	}
	if request.Description != nil {
		modifyCategoryRequest.Description = &wrappers.StringValue{Value: *request.Description}
	}
	if request.Icon != nil {
		modifyCategoryRequest.Icon = &wrappers.BytesValue{Value: request.Icon}
	}

	_, err = op.Category().ModifyCategory(openpitrix.SystemContext(), modifyCategoryRequest)
	if err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

func DescribeCategory(id string) (*Category, error) {
	op, err := cs.ClientSets().OpenPitrix()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	resp, err := op.Category().DescribeCategories(openpitrix.SystemContext(), &pb.DescribeCategoriesRequest{
		CategoryId: []string{id},
		Limit:      1,
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	var category *Category

	if len(resp.CategorySet) > 0 {
		category = convertCategory(resp.CategorySet[0])
		return category, nil
	} else {
		err := status.New(codes.NotFound, "resource not found").Err()
		klog.Error(err)
		return nil, err
	}
}

func ListCategories(conditions *params.Conditions, orderBy string, reverse bool, limit, offset int) (*models.PageableResponse, error) {
	client, err := cs.ClientSets().OpenPitrix()

	if err != nil {
		klog.Error(err)
		return nil, err
	}

	req := &pb.DescribeCategoriesRequest{}

	if keyword := conditions.Match["keyword"]; keyword != "" {
		req.SearchWord = &wrappers.StringValue{Value: keyword}
	}
	if orderBy != "" {
		req.SortKey = &wrappers.StringValue{Value: orderBy}
	}
	req.Reverse = &wrappers.BoolValue{Value: reverse}
	req.Limit = uint32(limit)
	req.Offset = uint32(offset)
	resp, err := client.Category().DescribeCategories(openpitrix.SystemContext(), req)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	items := make([]interface{}, 0)

	for _, item := range resp.CategorySet {
		items = append(items, convertCategory(item))
	}

	return &models.PageableResponse{Items: items, TotalCount: int(resp.TotalCount)}, nil
}

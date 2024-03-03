package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"

	"github.com/RustGrub/FunnyGoService/consts"
	"github.com/RustGrub/FunnyGoService/http/middleware"
	"github.com/RustGrub/FunnyGoService/http/utils"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"github.com/google/uuid"
)

// r.Context() в случае входящих запросов отменяется при потере соединения с клиентом, либо при return из ServeHttp

func (s *FunnyService) getGood(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	ctx := context.WithValue(r.Context(), consts.ReqID, reqID)

	goodID, err := middleware.IntValueFromRequest(r, "id")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	projectID, err := middleware.IntValueFromRequest(r, "projectId")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}
	res, err := s.usecases.GetGood(ctx, *goodID, *projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(err)
			resp := models.BadResponse{
				Code:    3,
				Message: consts.ErrGoodNotFound,
				Details: []string{"Some details..."},
			}
			if err = json.NewEncoder(w).Encode(&resp); err != nil {
				err = fmt.Errorf(consts.InternalServerError+": %v", err)
				s.logger.Error(err)
				utils.Write500(err.Error(), w)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = fmt.Errorf(consts.ErrInDatabase+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		err = fmt.Errorf(consts.InternalServerError+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}

}

func (s *FunnyService) createGood(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	ctx := context.WithValue(r.Context(), consts.ReqID, reqID)

	projectID, err := middleware.IntValueFromRequest(r, "projectId")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	var req models.CreateGoodRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}
	if middleware.IsEmptyOrOnlySpacesString(req.Name) {
		err = fmt.Errorf("empty name")
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	req.ProjectID = *projectID
	res, err := s.usecases.CreateGood(ctx, req)
	if err != nil {
		err = fmt.Errorf(consts.ErrInDatabase+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		err = fmt.Errorf(consts.InternalServerError+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
}

func (s *FunnyService) updateGood(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	ctx := context.WithValue(r.Context(), consts.ReqID, reqID)

	goodID, err := middleware.IntValueFromRequest(r, "id")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	projectID, err := middleware.IntValueFromRequest(r, "projectId")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	var req models.UpdateGoodRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}
	if middleware.IsEmptyOrOnlySpacesString(req.Name) {
		err = fmt.Errorf("empty name")
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}
	req.ProjectID = *projectID
	req.GoodID = *goodID

	res, err := s.usecases.UpdateGood(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(err)
			resp := models.BadResponse{
				Code:    3,
				Message: consts.ErrGoodNotFound,
				Details: []string{"Some details..."},
			}
			if err = json.NewEncoder(w).Encode(&resp); err != nil {
				err = fmt.Errorf(consts.InternalServerError+": %v", err)
				s.logger.Error(err)
				utils.Write500(err.Error(), w)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = fmt.Errorf(consts.ErrInDatabase+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		err = fmt.Errorf(consts.InternalServerError+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
}
func (s *FunnyService) removeGood(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	ctx := context.WithValue(r.Context(), consts.ReqID, reqID)

	goodID, err := middleware.IntValueFromRequest(r, "id")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	projectID, err := middleware.IntValueFromRequest(r, "projectId")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	req := models.RemoveGoodRequest{GoodID: *goodID, ProjectID: *projectID}

	res, err := s.usecases.RemoveGood(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(err)
			resp := models.BadResponse{
				Code:    3,
				Message: consts.ErrGoodNotFound,
				Details: []string{"Some details..."},
			}
			if err = json.NewEncoder(w).Encode(&resp); err != nil {
				err = fmt.Errorf(consts.InternalServerError+": %v", err)
				s.logger.Error(err)
				utils.Write500(err.Error(), w)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = fmt.Errorf(consts.ErrInDatabase+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		err = fmt.Errorf(consts.InternalServerError+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
}
func (s *FunnyService) getList(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	ctx := context.WithValue(r.Context(), consts.ReqID, reqID)

	limit, err := middleware.IntValueFromRequest(r, "limit")
	if err != nil {
		p := consts.DefaultLimit
		/*
			err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
			s.logger.Error(err)
			utils.Write400(err.Error(), w)
			return*/
		limit = &p
	}
	offset, err := middleware.IntValueFromRequest(r, "offset")
	if err != nil {
		p := consts.DefaultOffset
		/*
			err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
			s.logger.Error(err)
			utils.Write400(err.Error(), w)
			return*/
		offset = &p
	}

	res, err := s.usecases.GetGoodsList(ctx, *limit, *offset)
	if err != nil {
		err = fmt.Errorf(consts.ErrInDatabase+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		err = fmt.Errorf(consts.InternalServerError+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
}

func (s *FunnyService) reprioritizeGoods(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	ctx := context.WithValue(r.Context(), consts.ReqID, reqID)

	goodID, err := middleware.IntValueFromRequest(r, "id")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	projectID, err := middleware.IntValueFromRequest(r, "projectId")
	if err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	var req models.Reprioritize
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf(consts.ErrBadRequest+": %v", err)
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}
	if req.Priority <= 0 {
		err = fmt.Errorf(consts.ErrBadRequest + ": invalid priority")
		s.logger.Error(err)
		utils.Write400(err.Error(), w)
		return
	}

	req.ProjectID = *projectID
	req.GoodID = *goodID

	res, err := s.usecases.ReprioritizeGoods(ctx, req)
	if err != nil {
		err = fmt.Errorf(consts.ErrInDatabase+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		err = fmt.Errorf(consts.InternalServerError+": %v", err)
		s.logger.Error(err)
		utils.Write500(err.Error(), w)
		return
	}
}

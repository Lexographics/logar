package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sadk.dev/logar"
)

func (h *Handler) ListActions(w http.ResponseWriter, r *http.Request) {
	actionsMap := h.logger.GetActionManager().GetActionsMap()
	details := []ActionDetails{}

	for _, action := range actionsMap {
		argTypes, err := h.logger.GetActionManager().GetActionArgTypes(action.Path)
		if err != nil {
			h.logger.GetLogger().Error(logar.LogarLogs, fmt.Sprintf("Error getting arg types for action %s: %v", action.Path, err), "api")
			continue
		}

		argTypeData := make([]ArgType, len(argTypes))
		for i, t := range argTypes {
			kind, ok := h.logger.GetTypeKind(t)
			if !ok {
				kind = logar.TypeKind_Text
			}

			argTypeData[i] = ArgType{
				Type: t.String(),
				Kind: string(kind),
			}
		}

		details = append(details, ActionDetails{
			Path:        action.Path,
			Args:        argTypeData,
			Description: action.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, details))
}

func (h *Handler) InvokeActionHandler(w http.ResponseWriter, r *http.Request) {
	var req InvokeActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, fmt.Sprintf("Invalid request body: %v", err)))
		return
	}

	if req.Path == "" {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, "Missing 'path' in request body"))
		return
	}

	expectedTypes, err := h.logger.GetActionManager().GetActionArgTypes(req.Path)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, fmt.Sprintf("Error finding action '%s': %v", req.Path, err)))
		return
	}

	if len(req.Args) != len(expectedTypes) {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, fmt.Sprintf("Action '%s' expects %d arguments, but received %d", req.Path, len(expectedTypes), len(req.Args))))
		return
	}

	parsedArgs := make([]any, len(req.Args))
	for i, argStr := range req.Args {
		expectedType := expectedTypes[i]
		val, err := parseStringArg(argStr, expectedType)
		if err != nil {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(NewResponse(StatusCode_InvalidRequest, fmt.Sprintf("Error parsing argument %d for action '%s': expected type %s, error: %v", i+1, req.Path, expectedType.String(), err)))
			return
		}
		parsedArgs[i] = val
	}

	result, err := h.logger.GetActionManager().InvokeAction(req.Path, parsedArgs...)

	w.Header().Set("Content-Type", "application/json")
	resp := InvokeActionResponse{}
	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(NewResponse(StatusCode_Error, err.Error()))
	} else {
		if len(result) == 1 {
			resp.Result = result[0]
		} else {
			resp.Result = result
		}
	}

	json.NewEncoder(w).Encode(NewResponse(StatusCode_Success, resp))
}

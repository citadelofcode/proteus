package test

import (
	"github.com/citadelofcode/proteus/internal"
)

const (
	FOLDER_TYPE_PATH = internal.FOLDER_TYPE_PATH
	FILE_TYPE_PATH = internal.FILE_TYPE_PATH
	StatusOK = internal.StatusOK
)

var GetPathType = internal.GetPathType
var NewServer = internal.NewServer
var NewRouter = internal.NewRouter
var NormalizeRoute = internal.NormalizeRoute
var EmptyPrefixTree = internal.EmptyPrefixTree
var GetHighestVersion = internal.GetHighestVersion
var GetResponseVersion = internal.GetResponseVersion

type FileSystemError = internal.FileSystemError
type ResponseError = internal.ResponseError
type RoutingError = internal.RoutingError
type Headers = internal.Headers
type Params = internal.Params
type HttpRequest = internal.HttpRequest
type HttpResponse = internal.HttpResponse
type StatusCode = internal.StatusCode
type Route = internal.Route

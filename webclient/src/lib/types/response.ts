
export type StatusCode = number;

export const StatusCode = {
  // Unknown values
	Unknown: 0,

  // Success values
	Success: 1,

  // Error values
	Error: 1000,
	SessionExpired: 1001,
	InvalidRequest: 1002,
	InvalidCredentials: 1003,
}

export type Response<T> = {
	status_code: StatusCode;
	data: T;
}

export function isSuccess(statusCode: StatusCode): boolean {
	return statusCode >= StatusCode.Success && statusCode < StatusCode.Error;
}

export function isError(statusCode: StatusCode): boolean {
	return statusCode >= StatusCode.Error || statusCode === StatusCode.Unknown;
}


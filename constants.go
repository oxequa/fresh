package fresh

// MIME types
const (
	MIMEAppJSON    = "application/json" + ";" + UTF8
	MIMEAppJS      = "application/javascript" + ";" + UTF8
	MIMEAppXML     = "application/xml" + ";" + UTF8
	MIMEUrlencoded = "application/x-www-form-urlencoded"
	MIMEMultipart  = "multipart/form-data"
	MIMETextHTML   = "text/html" + ";" + UTF8
	MIMETextXML    = "text/xml" + ";" + UTF8
	MIMEText       = "text/plain" + ";" + UTF8
	MIMEGzip       = "gzip"
)

// Access
const (
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllows           = "Access-Control-Allow-s"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlExposes          = "Access-Control-Expose-s"
	AccessControlRequests         = "Access-Control-Request-s"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// Request
const (
	Accept              = "Accept"
	AcceptEncoding      = "Accept-Encoding"
	Allow               = "Allow"
	Authorization       = "Authorization"
	ContentDisposition  = "Content-Disposition"
	ContentEncoding     = "Content-Encoding"
	ContentLength       = "Content-Length"
	ContentType         = "Content-Type"
	Cookie              = "Cookie"
	SetCookie           = "Set-Cookie"
	IfModifiedSince     = "If-Modified-Since"
	LastModified        = "Last-Modified"
	Location            = "Location"
	Upgrade             = "Upgrade"
	Vary                = "Vary"
	WWWAuthenticate     = "WWW-Authenticate"
	XForwardedFor       = "X-Forwarded-For"
	XForwardedProto     = "X-Forwarded-Proto"
	XForwardedProtocol  = "X-Forwarded-Protocol"
	XForwardedSsl       = "X-Forwarded-Ssl"
	XUrlScheme          = "X-Url-Scheme"
	XHTTPMethodOverride = "X-HTTP-Method-Override"
	XRealIP             = "X-Real-IP"
	XRequestID          = "X-Request-ID"
	Server              = "Server"
	Origin              = "Origin"
)

// Encoding chartset
const (
	UTF8     = "charset=UTF-8"
	ISO88591 = "chartset=ISO-8859-1"
)

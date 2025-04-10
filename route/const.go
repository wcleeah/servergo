package route

var codeMsgMap = map[string]string{
	// 1xx Informational Responses
	"100": "Continue",
	"101": "Switching Protocols",
	"102": "Processing",
	"103": "Early Hints",

	// 2xx Successful Responses
	"200": "OK",
	"201": "Created",
	"202": "Accepted",
	"203": "Non-Authoritative Information",
	"204": "No Content",
	"205": "Reset Content",
	"206": "Partial Content",
	"207": "Multi-Status",
	"208": "Already Reported",
	"226": "IM Used",

	// 3xx Redirection Messages
	"300": "Multiple Choices",
	"301": "Moved Permanently",
	"302": "Found",
	"303": "See Other",
	"304": "Not Modified",
	"305": "Use Proxy",    // Deprecated
	"306": "Switch Proxy", // Reserved but unused
	"307": "Temporary Redirect",
	"308": "Permanent Redirect",

	// 4xx Client Error Responses
	"400": "Bad Request",
	"401": "Unauthorized",
	"402": "Payment Required", // Rarely used
	"403": "Forbidden",
	"404": "Not Found",
	"405": "Method Not Allowed",
	"406": "Not Acceptable",
	"407": "Proxy Authentication Required",
	"408": "Request Timeout",
	"409": "Conflict",
	"410": "Gone",
	"411": "Length Required",
	"412": "Precondition Failed",
	"413": "Payload Too Large", // Previously Content Too Large
	"414": "URI Too Long",
	"415": "Unsupported Media Type",
	"416": "Range Not Satisfiable",
	"417": "Expectation Failed",
	"418": "I'm a Teapot", // Joke status code from RFC2324
	"421": "Misdirected Request",
	"422": "Unprocessable Entity",
	"423": "Locked",
	"424": "Failed Dependency",
	"425": "Too Early", // Experimental
	"426": "Upgrade Required",
	"428": "Precondition Required",
	"429": "Too Many Requests",
	"431": "Request Header Fields Too Large",
	"451": "Unavailable For Legal Reasons",

	// 5xx Server Error Responses
	"500": "Internal Server Error",
	"501": "Not Implemented",
	"502": "Bad Gateway",
	"503": "Service Unavailable",
	"504": "Gateway Timeout",
	"505": "HTTP Version Not Supported",
	"506": "Variant Also Negotiates",         // Rarely used
	"507": "Insufficient Storage",            // WebDAV specific
	"508": "Loop Detected",                   // WebDAV specific
	"510": "Not Extended",                    // RFC2774 specific
	"511": "Network Authentication Required", // Captive portals
}

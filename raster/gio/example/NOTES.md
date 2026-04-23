An .ivg file does not preserve the width and height attributes of the original SVG file.
Additionally, the preserveAspectRatio attribute is also ignored.

According to the SVG specification, default value when the preserveAspectRatio attribute
is not specified is "xMidYMid meet". This means that the image is scaled to fit the viewport
while preserving the aspect ratio. The image is centered in the viewport along the x and y axes.

type PreserveAspectRatio int

const (
	// XMidYMidMeet is the default value of the PreserveAspectRatio attribute.
	XMidYMidMeet PreserveAspectRatio = iota
	XMinYMinMeet
	XMaxYMaxMeet
	XMidYMinMeet
	XMidYMaxMeet
	XMinYMidMeet
	XMaxYMidMeet
	XMinYMaxMeet
	XMaxYMinMeet
	XMidYMidSlice
	XMinYMinSlice
	XMaxYMaxSlice
	XMidYMinSlice
	XMidYMaxSlice
	XMinYMidSlice
	XMaxYMidSlice
	XMinYMaxSlice
	XMaxYMinSlice
	None
)

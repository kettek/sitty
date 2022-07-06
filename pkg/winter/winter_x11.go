//go:build linux

package winter

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Winter struct {
	ctx  *xgb.Conn
	root xproto.Window
}

func NewWinter() (*Winter, error) {
	w := &Winter{}
	c, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	w.ctx = c

	// Get the window id of the root window.
	setup := xproto.Setup(w.ctx)

	w.root = setup.DefaultScreen(w.ctx).Root

	return w, err
}

func (w *Winter) GetActiveWindowDimensions() (x, y, width, height int, err error) {
	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(w.ctx, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		log.Fatal(err)
	}

	// Get the atom id (i.e., intern an atom) of "_NET_WM_NAME".
	aname = "_NET_WM_NAME"
	nameAtom, err := xproto.InternAtom(w.ctx, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		log.Fatal(err)
	}

	reply, err := xproto.GetProperty(w.ctx, false, w.root, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Fatal(err)
	}
	windowId := xproto.Window(xgb.Get32(reply.Value))

	reply, err = xproto.GetProperty(w.ctx, false, windowId, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(reply.Value))
	if string(reply.Value) == "Go Sitter" {
		return 0, 0, 0, 0, fmt.Errorf("self")
	}

	// Let's get our window border size
	aname = "_NET_FRAME_EXTENTS"
	nameAtom, err = xproto.InternAtom(w.ctx, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		log.Fatal(err)
	}

	reply, err = xproto.GetProperty(w.ctx, false, windowId, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Fatal(err)
	}
	if len(reply.Value) >= 16 {
		left := binary.LittleEndian.Uint32(reply.Value[0:4])
		//right := binary.LittleEndian.Uint32(reply.Value[4:8])
		top := binary.LittleEndian.Uint32(reply.Value[8:12])
		//bottom := binary.LittleEndian.Uint32(reply.Value[12:16])
		x += int(left)
		y -= int(top)
	}

	gc := xproto.GetGeometry(w.ctx, xproto.Drawable(windowId))

	greply, gerr := gc.Reply()
	if gerr != nil {
		log.Fatal(gerr)
	}
	width += int(greply.Width)
	height += int(greply.Height)

	cc := xproto.TranslateCoordinates(w.ctx, windowId, w.root, 0, 0)

	creply, cerr := cc.Reply()
	if cerr != nil {
		log.Fatal(cerr)
	}

	// Check if we'd be out of bounds in useable work area.
	aname = "_NET_WORKAREA"
	nameAtom, err = xproto.InternAtom(w.ctx, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		log.Fatal(err)
	}
	reply, err = xproto.GetProperty(w.ctx, false, w.root, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Somehow use workarea

	x += int(creply.DstX)
	y += int(creply.DstY)

	x += int(greply.BorderWidth)
	y += int(greply.BorderWidth)

	return
}

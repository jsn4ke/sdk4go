package raft_ext

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func Gen(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + ".ext.go"

	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P(`package `, file.GoPackageName)
	head := `
	import (
		jsn_rpc "github.com/jsn4ke/jsn_net/rpc"
		"google.golang.org/protobuf/proto"
	)`
	g.P(head)

	needKuohao := false
	for _, v := range file.Messages {
		if !strings.HasSuffix(v.GoIdent.GoName, `Request`) &&
			!strings.HasSuffix(v.GoIdent.GoName, `Response`) {
			continue
		}
		if !needKuohao {
			needKuohao = true
			g.P(`var (`)

		}
		g.P(`_ jsn_rpc.RpcUnit = (*`, v.GoIdent.GoName, `)(nil)`)
	}
	if needKuohao {
		g.P(`)`)
	}
	for _, v := range file.Messages {
		if !strings.HasSuffix(v.GoIdent.GoName, `Request`) &&
			!strings.HasSuffix(v.GoIdent.GoName, `Response`) {
			continue
		}
		g.P(`func (*`, v.GoIdent.GoName, `) CmdId() uint32 {`)
		g.P(`return uint32(Cmd_Cmd_`, v.GoIdent.GoName, `)`)
		g.P(`}`)

		g.P(`func (x *`, v.GoIdent.GoName, `) Marshal() ([]byte, error) {`)
		g.P(`return proto.Marshal(x)`)
		g.P("}")

		g.P(`func (x *`, v.GoIdent.GoName, `) Unmarshal(in []byte)  error {`)
		g.P(`return proto.Unmarshal(in, x)`)
		g.P("}")

		g.P(`func (*`, v.GoIdent.GoName, `) New() jsn_rpc.RpcUnit {`)
		g.P(`return new(`, v.GoIdent.GoName, `)`)
		g.P(`}`)
	}
}

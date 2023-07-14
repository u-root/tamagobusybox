package bb

var bbRegisterSource = []byte("// Copyright 2018 the u-root Authors. All rights reserved\n// Use of this source code is governed by a BSD-style\n// license that can be found in the LICENSE file.\n\n// Package bbmain is the package imported by all rewritten busybox\n// command-packages to register themselves.\npackage bbmain\n\nimport (\n\t\"errors\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t// There MUST NOT be any other dependencies here.\n\t//\n\t// It is preferred to copy minimal code necessary into this file, as\n\t// dependency management for this main file is... hard.\n)\n\n// ErrNotRegistered is returned by Run if the given command is not registered.\nvar ErrNotRegistered = errors.New(\"command is not present in busybox\")\n\n// Noop is a noop function.\nvar Noop = func() {}\n\n// ListCmds lists bb commands and verifies symlinks.\n// It is by convention called when the bb command is invoked directly.\n// For every command, there should be a symlink in /bbin,\n// and for every symlink, there should be a command.\n// Occasionally, we have bugs that result in one of these\n// being false. Just running bb is an easy way to tell if something\n// in your image is messed up.\nfunc ListCmds() {\n\ttype known struct {\n\t\tname string\n\t\tbb   string\n\t}\n\tnames := map[string]*known{}\n\tg, err := filepath.Glob(\"/bbin/*\")\n\tif err != nil {\n\t\tfmt.Printf(\"bb: unable to enumerate /bbin\")\n\t}\n\n\t// First step is to assemble a list of all possible\n\t// names, both from /bbin/* and our built in commands.\n\tfor _, l := range g {\n\t\tif l == \"/bbin/bb\" {\n\t\t\tcontinue\n\t\t}\n\t\tb := filepath.Base(l)\n\t\tnames[b] = &known{name: l}\n\t}\n\tfor n := range bbCmds {\n\t\tif n == \"bb\" {\n\t\t\tcontinue\n\t\t}\n\t\tif c, ok := names[n]; ok {\n\t\t\tc.bb = n\n\t\t\tcontinue\n\t\t}\n\t\tnames[n] = &known{bb: n}\n\t}\n\t// Now walk the array of structs.\n\t// We don't sort as we don't want the\n\t// footprint of bringing in the package.\n\t// If you want it sorted, bb | sort\n\tvar hadError bool\n\tfor c, k := range names {\n\t\tif len(k.name) == 0 || len(k.bb) == 0 {\n\t\t\thadError = true\n\t\t\tfmt.Printf(\"%s:\\t\", c)\n\t\t\tif k.name == \"\" {\n\t\t\t\tfmt.Printf(\"NO SYMLINK\\t\")\n\t\t\t} else {\n\t\t\t\tfmt.Printf(\"%q\\t\", k.name)\n\t\t\t}\n\t\t\tif k.bb == \"\" {\n\t\t\t\tfmt.Printf(\"NO COMMAND\\n\")\n\t\t\t} else {\n\t\t\t\tfmt.Printf(\"%s\\n\", k.bb)\n\t\t\t}\n\t\t}\n\t}\n\tif hadError {\n\t\tfmt.Println(\"There is at least one problem. Known causes:\")\n\t\tfmt.Println(\"At least two initrds -- one compiled in to the kernel, a second supplied by the bootloader.\")\n\t\tfmt.Println(\"The initrd cpio was changed after creation or merged with another one.\")\n\t\tfmt.Println(\"When the initrd was created, files were inserted into /bbin by mistake.\")\n\t\tfmt.Println(\"Post boot, files were added to /bbin.\")\n\t}\n}\n\ntype bbCmd struct {\n\tinit, main func()\n}\n\nvar bbCmds = map[string]bbCmd{}\n\nvar defaultCmd *bbCmd\n\n// Register registers an init and main function for name.\nfunc Register(name string, init, main func()) {\n\tif _, ok := bbCmds[name]; ok {\n\t\tpanic(fmt.Sprintf(\"cannot register two commands with name %q\", name))\n\t}\n\tbbCmds[name] = bbCmd{\n\t\tinit: init,\n\t\tmain: main,\n\t}\n}\n\n// RegisterDefault registers a default init and main function.\nfunc RegisterDefault(init, main func()) {\n\tdefaultCmd = &bbCmd{\n\t\tinit: init,\n\t\tmain: main,\n\t}\n}\n\n// Run runs the command with the given name.\n//\n// If the command's main exits without calling os.Exit, Run will exit with exit\n// code 0.\nfunc Run(name string) error {\n\tvar cmd *bbCmd\n\tif c, ok := bbCmds[name]; ok {\n\t\tcmd = &c\n\t} else if defaultCmd != nil {\n\t\tcmd = defaultCmd\n\t} else {\n\t\treturn fmt.Errorf(\"%w: %s\", ErrNotRegistered, name)\n\t}\n\tcmd.init()\n\tcmd.main()\n\tos.Exit(0)\n\t// Unreachable.\n\treturn nil\n}\n")

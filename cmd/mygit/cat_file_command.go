package main

import (
	"fmt"
	"os"
	"io"
	"compress/zlib"
	"path/filepath"
	"strconv"
	"bufio"
)

func catFileCmd(args []string) (err error) {
	if len(args) < 3 || args[1] != "-p" {
		fmt.Fprint(os.Stderr, "usage: mygit cat-file -p <blob_hash>\n")

		return fmt.Errorf("bad bad")
	}

	sha := args[2]

	path := filepath.Join(".git", "objects", sha[:2], sha[2:])
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	defer func() {
		e := file.Close()
		if err == nil && e != nil { 
			err = fmt.Errorf("%w", e)
		}
	}()
	return catFile(file)
}


func catFile(r io.Reader) (err error) {
	zr, err := zlib.NewReader(r)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		e := zr.Close()
		if err == nil && e != nil { 
			err = fmt.Errorf("%w", e)
		}
	}()
	err = parseObject(zr)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func parseObject(r io.Reader) (err error) {
	br := bufio.NewReader(r)

	typ, err := br.ReadString(' ')
	if err != nil {
		return err
	}

	typ = typ[:len(typ)-1]

	if typ != "blob" {
		return fmt.Errorf("unsupported type: %v", typ)
	}

	sizeStr, err := br.ReadString('\000')
	if err != nil {
		return err
	}
	sizeStr = sizeStr[:len(sizeStr)-1]
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse size: %w", err)
	}

	_, err = io.CopyN(os.Stdout, br, size)
	if err != nil {
		return err
	}
	return nil
}

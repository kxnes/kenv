// Package parser contains instruments for parsing AST and target.
package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/kxnes/kenv/internal/conv"
	"github.com/kxnes/kenv/internal/types"
)

// ErrParse decorates the list of parsing errors.
// Must contain only ErrSyntax and ErrConsistence errors
type ErrParse []error

func (errs ErrParse) Error() string {
	var s []string
	for num, err := range errs {
		s = append(s, fmt.Sprintf("%02d. %s", num+1, err))
	}
	return fmt.Sprintf("parse errors:\n%s\n", strings.Join(s, "\n"))
}

// ErrSyntax occurs while parsing target structure fields.
type ErrSyntax struct {
	line int
	errs []error
}

func (e *ErrSyntax) Error() string {
	var s []string
	for _, err := range e.errs {
		s = append(s, fmt.Sprintf("%q", err))
	}
	return fmt.Sprintf("[syntax] line %d errors: %s", e.line, strings.Join(s, ", "))
}

// append appends error `err` and re-translates `val`.
func (e *ErrSyntax) append(val interface{}, err error) interface{} {
	if err != nil {
		e.errs = append(e.errs, err)
	}
	return val
}

// ErrConsistence occurs if convert function was not found for field type.
type ErrConsistence string

func (e ErrConsistence) Error() string {
	return fmt.Sprintf("[consistance] %s", string(e))
}

// Parser uses for parsing incoming file.
type Parser struct {
	// Name of search target structure.
	target string
	// Name of file that contains target.
	filename string
	// All syntax and consistence errors.
	errors ErrParse
	// Parsed target struct fields that will contain environment variables.
	fields map[string]*types.Field
	// Parsed target struct methods that will uses for type conversion.
	typeConv map[string]string
	// AST representation of file entities.
	asTree *ast.File
	// Inner file representation for showing line numbers.
	file *token.File
}

// New returns new target parser.
func New(target, filename string) *Parser {
	return &Parser{
		target:   target,
		filename: filename,
		fields:   make(map[string]*types.Field),
		typeConv: make(map[string]string),
	}
}

// Parse parses file into inner `Parser.fields` and `Parser.typeConv`.
// Also after parsing it is guarantee that all type conversion functions exist.
func (p *Parser) Parse() (err error) {
	// for safe method usage
	if p == nil {
		return errors.New("parser not initialized")
	}

	fs := token.NewFileSet()
	p.asTree, err = parser.ParseFile(fs, p.filename, nil, 0)
	if err != nil {
		return fmt.Errorf("parse file error: %w", err)
	}

	// because we have only one file and use start position
	p.file = fs.File(1)

	ast.Walk(p, p.asTree) // syntax errors

	if len(p.fields) == 0 {
		return fmt.Errorf("target %q not found in file %q", p.target, p.filename)
	}

	p.checkConsistency() // consistence errors

	if len(p.errors) != 0 {
		return p.errors
	}

	return
}

// checkConsistency checks that all found environment variables have convert functions.
// If not user-defined will be use predefined.
func (p *Parser) checkConsistency() {
	predefined := conv.GetPredefinedConv()

	for name, t := range p.fields {
		if fn, ok := p.typeConv[t.Type]; ok {
			t.Func = fn
			continue
		}

		if pre, ok := predefined[t.Type]; ok {
			t.Func = pre
			continue
		}

		err := ErrConsistence(fmt.Sprintf("no converter for %q (%q)", t.Type, name))
		p.errors = append(p.errors, err)
	}
}

func (p *Parser) Visit(node ast.Node) ast.Visitor {
	// see the `ast.Visitor`'s `Visit` method docs
	if node == nil {
		return nil
	}

	switch decl := node.(type) {
	case *ast.GenDecl:
		p.findEnvVars(decl)
	case *ast.FuncDecl:
		p.findTypeConv(decl)
	}

	return p
}

// findEnvVars searches for environment variable of target struct.
func (p *Parser) findEnvVars(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			break
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok || ts.Name.Name != p.target {
			continue
		}

		// because names unique we can omit checks like "Did we find already?"
		p.parse("", st)
		return
	}
}

// parse does inner parsing for struct type.
func (p *Parser) parse(prefix string, s *ast.StructType) {
	// stop recursion for nested structs
	if s == nil {
		return
	}

	for _, f := range s.Fields.List {
		err := &ErrSyntax{line: p.file.Line(f.Pos())}

		var (
			tag    = err.append(p.parseTag(f)).(*types.Tag)
			name   = prefix + err.append(p.parseName(f)).(string)
			typeID = err.append(p.parseType(f)).(string)
		)

		switch {
		case err.errs != nil:
			p.errors = append(p.errors, err)
		case tag == nil:
			p.parse(name+".", p.parseSubType(f))
		default:
			if tag.EnvVar == "" {
				tag.EnvVar = strings.ToUpper(strings.ReplaceAll(name, ".", "_"))
			}

			// because names unique (because of prefix) we can omit checks like "Did we find already?"
			p.fields[name] = &types.Field{
				Name:    name,
				Type:    typeID,
				EnvVar:  tag.EnvVar,
				Action:  tag.Action,
			}
		}
	}
}

// parseTag parses the tag of field `f`.
func (p *Parser) parseTag(f *ast.Field) (*types.Tag, error) {
	// tag not exist and this is not error (maybe inner type will contain it)
	if f.Tag == nil {
		return nil, nil
	}

	var targetTag string

	tags := strings.Trim(f.Tag.Value, "`")
	for _, tag := range strings.Fields(tags) {
		if strings.HasPrefix(tag, "env:") {
			targetTag = tag
			break
		}
	}

	if targetTag == "" {
		return nil, errors.New("tag not found")
	}

	rawBody := strings.Split(targetTag, `"`)
	if len(rawBody) != 3 {
		return nil, errors.New("invalid tag")
	}

	tagBody := rawBody[1]
	envVar := tagBody[:len(tagBody)-1]
	action := tagBody[len(tagBody)-1:]

	switch action {
	case types.Want, types.Must, types.Secret:
		return &types.Tag{EnvVar: envVar, Action: action}, nil
	default:
		return nil, errors.New("invalid tag format")
	}
}

// parseName parses the name of field `f`.
func (p *Parser) parseName(f *ast.Field) (string, error) {
	if len(f.Names) > 1 {
		return "", errors.New("multiple names")
	}

	if len(f.Names) == 0 {
		return "", errors.New("embedding not supported")
	}

	name := f.Names[0]
	if !name.IsExported() {
		return "", errors.New("field not exported")
	}

	// name must exist
	return name.Name, nil
}

// parseType parses the type of field `f`.
func (p *Parser) parseType(f *ast.Field) (string, error) {
	var typeID string

	switch t := f.Type.(type) {
	case *ast.Ident:
		typeID = t.Name
	case *ast.SelectorExpr:
		typeID = t.X.(*ast.Ident).Name + "." + t.Sel.Name
	case *ast.StructType:
		return "", errors.New("unnamed type not supported")
	default:
		return "", errors.New("type not supported")
	}

	return typeID, nil
}

// parseSubType parses the subtype of field `f`. Subtype may be only struct.
func (p *Parser) parseSubType(f *ast.Field) *ast.StructType {
	t, ok := f.Type.(*ast.Ident)
	if !ok {
		return nil
	}

	if t.Obj == nil {
		return nil
	}

	spec, ok := t.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	subType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	return subType
}

// findTypeConv searches convert functions for environment variable of target struct.
func (p *Parser) findTypeConv(decl *ast.FuncDecl) {
	// check receiver signature
	if decl.Recv == nil {
		return
	}

	var recvType string

	recv := decl.Recv.List[0] // guarantee only one receiver
	switch t := recv.Type.(type) {
	case *ast.Ident:
		recvType = t.Name
	case *ast.StarExpr:
		recvType = t.X.(*ast.Ident).Name
	}

	if recvType != p.target {
		return
	}

	// check params signature
	params := decl.Type.Params
	if params == nil || len(params.List) != 1 {
		return
	}

	paramType := params.List[0].Type
	if paramType.(*ast.Ident).Name != "string" {
		return
	}

	// check return signature
	results := decl.Type.Results
	if results == nil || len(results.List) != 2 {
		return
	}

	resultType, errorType := results.List[0].Type, results.List[1].Type
	if errorType.(*ast.Ident).Name != "error" {
		return
	}

	var targetType string
	switch t := resultType.(type) {
	case *ast.Ident:
		targetType = t.Name
	case *ast.SelectorExpr:
		targetType = t.X.(*ast.Ident).Name + "." + t.Sel.Name
	default:
		return
	}

	p.typeConv[targetType] = decl.Name.Name
}

// Overview prints parsing results as a table.
func (p *Parser) Overview() {
	format := "%20s | %20s | %20s | %6s | %20s\n"
	fmt.Printf(format, "Name", "Type", "Env Variable", "Action", "Converter")
	fmt.Println(strings.Repeat("-", 20+3+20+3+20+3+6+3+20))
	for n, v := range p.fields {
		fmt.Printf(format, n, v.Type, v.EnvVar, v.Action, v.Func)
	}
}

// ParsedFields returns parsed fields only if parser parsed.
func (p *Parser) ParsedFields() map[string]*types.Field {
	if len(p.errors) != 0 || len(p.fields) == 0 {
		return nil
	}
	return p.fields
}

// Package returns name of package that contain target.
func (p *Parser) Package() string {
	return p.asTree.Name.Name
}

// Target returns target name.
func (p *Parser) Target() string {
	return p.target
}

// Filename returns name of parsed file.
func (p *Parser) Filename() string {
	return p.file.Name()
}

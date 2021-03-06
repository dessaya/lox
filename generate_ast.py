import os, sys
import subprocess

args = sys.argv[1:]

if len(args) != 1:
      print("Usage: generate_ast <output directory>")
      sys.exit(64)

outputDir = args[0]

def defineAst(outputDir, baseName, types):
    path = outputDir + "/" + baseName.lower() + ".go"
    with open(path, "w") as f:
        f.write(f"""
            package lox

            // autogenerated with `python {' '.join(sys.argv)}`

            type {baseName} interface{{
                accept(v {baseName}Visitor) interface{{}}
            }}

        """)

        # The AST classes.
        for type in types:
            className = type.split(":")[0].strip()
            fields = type.split(":")[1].strip()
            defineType(f, baseName, className, fields)
        defineVisitor(f, baseName, types)
    subprocess.run(["gofmt", "-w", path])

def defineVisitor(f, baseName, types):
    f.write(f"type {baseName}Visitor interface {{\n");
    for type in types:
      typeName = type.split(":")[0].strip()
      f.write(f"visit{typeName}{baseName}({typeName[0].lower()} *{typeName}) interface{{}}\n")
    f.write("  }\n\n")



def defineType(f, baseName, className, fieldList):
    fields = fieldList.split(", ")

    f.write(f"type {className} struct {{\n")

    # Fields.
    for field in fields:
      f.write(f"    {field}\n")

    f.write("}\n")

    # Constructor.
    f.write(f"func New{className}({fieldList}) *{className} {{\n")
    f.write(f"return &{className} {{\n")

    # Store parameters in fields.
    for field in fields:
      name = field.split(" ")[0]
      f.write(f"{name}: {name},\n")

    f.write("    }\n")
    f.write("    }\n")

    f.write(f"""
        func ({className[0].lower()} *{className}) accept({baseName[0].lower()}v {baseName}Visitor) interface{{}} {{
            return {baseName[0].lower()}v.visit{className}{baseName}({className[0].lower()})
        }}
    """)

defineAst(outputDir, "Expr", [
    "Assign   : name *Token, value Expr",
    "Binary   : left Expr, operator *Token, right Expr",
    "Call     : callee Expr, paren *Token, arguments []Expr",
    "Grouping : expression Expr",
    "Literal  : value interface{}",
    "Logical  : left Expr, operator *Token, right Expr",
    "Unary    : operator *Token, right Expr",
    "Variable : name *Token",
])

defineAst(outputDir, "Stmt", [
    "Block      : statements []Stmt",
    "Expression : expression Expr",
    "Function   : name *Token, params []*Token, body []Stmt",
    "If         : condition Expr, thenBranch Stmt, elseBranch Stmt",
    "Print      : expression Expr",
    "Return     : keyword *Token, value Expr",
    "Var        : name *Token, initializer Expr",
    "While      : condition Expr, body Stmt",
]);

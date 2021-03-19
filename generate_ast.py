import os, sys
import subprocess

args = sys.argv[1:]

if len(args) != 1:
      print("Usage: generate_ast <output directory>")
      os.exit(64)

outputDir = args[0]

def defineAst(outputDir, baseName, types):
    path = outputDir + "/" + baseName.lower() + ".go"
    with open(path, "w") as f:
        f.write(f"""
            package main

            type {baseName} interface{{
                accept(v {baseName}Visitor)
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
      f.write(f"visit{typeName}({typeName[0].lower()} *{typeName})\n")
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
        func ({className[0].lower()} *{className}) accept({baseName[0].lower()}v {baseName}Visitor) {{
            ev.visit{className}({className[0].lower()})
        }}
    """)

defineAst(outputDir, "Expr", [
    "Binary   : left Expr, operator *Token, right Expr",
    "Grouping : expression Expr",
    "Literal  : value interface{}",
    "Unary    : operator *Token, right Expr",
])

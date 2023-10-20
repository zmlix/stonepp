elements  : expr { "," expr }
param     : Identifier
params    : param { "," param }
param_list: "(" [ params ] ")"
def       : "def" Identifier param_list block
member    : def | simple
class_body: "{" [ member ] { EOL [ member ] } "}"
defclass  : "class" Identifier [ "extends" Identifier ] class_body
args      : expr { "," expr }
postfix   : "." Identifier | "(" [ args ] ")" | "[" expr "]"
primary   : ("fun" param_list block | "[" [ elements ] "]" | "(" expr ")" | Number | Identifier | String | Boolean) { postfix }
factor    : {"-"} primary
expr      : factor { Op factor}
block     : "{" [ statement ] { EOL [ statement ] } "}"
simple    : expr [ args ] | "return" [ expr ]
statement : "if" expr block { "elif" expr block } [ "else" block ]
          | "while" expr block
          | simple
program   : [ defclass | def | statement ] EOL

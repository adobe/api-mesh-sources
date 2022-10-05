const root = process.argv[2]
const paths = process.argv.slice(3);
const fs = require('fs');
const JsonInterpolate = require('json-interpolate');

for (path of paths) {
    const source = JSON.parse(fs.readFileSync(`${root}/${path}`).toString());
    const jsonInterpolate = new JsonInterpolate({variablesSchema: source.variables});
    const variables = jsonInterpolate.getJsonVariables(JSON.stringify(source.provider));
    const notDeclared = jsonInterpolate.getVariablesWithoutDeclaredInterface(variables);

    if (notDeclared.length) {
        throw new Error(`
        \nThe schema for next variables is not declared: 
        ${notDeclared.map((element) => '\n' + element.name)}
        `)
    }
}


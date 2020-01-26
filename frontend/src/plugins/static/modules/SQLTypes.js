
var types = {

    postgres: {
        NUMERIC: genreateObject([

            'SMALLINT',
            'INTEGER',
            'NUMBER',
            'BIGINT',
            'DECIMAL',
            'NUMERIC',
            'REAL',
            'DOUBLE PRECISION',
            'SMALLSERIAL',
            'SERIAL',
            'BIGSERIAL',
            'MONEY'
        ]),
        STRING: genreateObject([
            'VARYING',
            'VARCHAR',
            'CHARACTER',
            'CHAR',
            'TEXT'
        ]),
        BINARY: genreateObject([
            'BYTEA'

        ]),
        DATE: genreateObject([
            'TIMESTAMP',
            'TIMESTAMPZ',
            'DATE',
            'TIME',
            'INTERVAL'
        ]),
        BOOLEAN: genreateObject([
            'BOOLEAN'
        ]),
        GEOMETRIC: genreateObject([
            'POINT',
            'LINE',
            'LSEG',
            'BOX',
            'PATH',
            'POLYGON',
            'CIRCLE'
        ])
    }

}

function genreateObject (array) {
    let i
    const obj = {}
    for (i in array) {
        obj[array[i].toUpperCase()] = array[i].toUpperCase()
    }

    obj.values = array

    return obj
}

export default types

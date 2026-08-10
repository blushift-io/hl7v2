var delimiters = {
    fieldSeparator: "|",
    componentSeparator: "^",
    subcomponentSeparator: "&",
    escapeCharacter: "\\",
    repetitionCharacter: "~",
    segmentSeparator: '\r'
}

function component() {
    this.value = []

    if (arguments.length > 0) {
        for (var i = 0; i < arguments.length; i++) {
            this.value.push(arguments[i])
        }
    }
}

component.prototype.toString = function (delimiters) {
    var res = "";

    for (var i = 0; i < this.value.length; i++) {
        if (Array.isArray(this.value[i])) {
            for (var j = 0; j < this.value[i].length; j++) {
                res += this.value[i][j];
                if (j != this.value[i].length - 1) res += delimiters.subcomponentSeparator
            }
        } else {
            res += this.value[i];
        }

        if (i != this.value.length - 1) res += delimiters.repetitionCharacter
    }

    return res;
}

function field() {
    this.value = [];

    if (arguments.length > 0) {
        for (var i = 0; i < arguments.length; i++) {
            if (Array.isArray(arguments[i])) {
                var components = new Array();
                for (var j = 0; j < arguments[i].length; j++) {
                    components.push(new component(arguments[i][j]))
                }

                this.value.push(components)
            } else {
                this.value.push(new component(arguments[i]));
            }
        }
    }
}

field.prototype.toString = function (delims) {
    if (!delims) {
        delims = delimiters
    }
    var res = "";

    for (var i = 0; i < this.value.length; i++) {
        if (Array.isArray(this.value[i])) {
            for (var j = 0; j < this.value[i].length; j++) {
                res += this.value[i][j].toString(delims);
                if (j != this.value[i].length - 1) res += delims.componentSeparator
            }
        } else {
            res += this.value[i];
        }

        if (i != this.value.length - 1) res += delims.repetitionCharacter
    }

    return res;
}

var segment = function () {
    this.name = ""
    this.fields = [];

    if (arguments.length >= 1) {
        this.name = arguments[0]
    }

    if (arguments.length >= 2) {
        for (var i = 1; i < arguments.length; i++) {
            if (Array.isArray(arguments[i])) {
                var fields = new Array();
                for (var j = 0; j < arguments[i].length; j++) {
                    fields.push(new field(arguments[i][j]))
                }

                this.fields.push(fields)
            } else {
                for (var i = 0; i < arguments.length; i++) {
                    this.fields.push(new field(arguments[i]))
                }
            }
        }
    }
}

segment.prototype.addField = function (val, pos) {
    if (pos) {
        if (this.fields.length > (pos - 1)) {
            this.setField(pos, val)
        } else {
            curLen = this.fields.length;
            while (curLen <= (pos - 2)) {
                this.addField("");
                curLen = this.fields.length;
            }

            this.addField(val)
        }
    } else {
        if ((typeof val) == 'object') {
            if (Array.isArray(val)) {
                this.fields.push(new field(val))
            } else {
                this.fields.push(val)
            }
        } else {
            this.fields.push(new field(val))
        }
    }
}

segment.prototype.setField = function (idx, val) {
    if (this.fields.length >= idx) {
        this.fields[idx - 1] = new field(val)
    }
}

segment.prototype.removeField = function (idx) {
    if (this.fields.length >= idx) {
        this.fields.splice(idx - 1, 1)
    }
}

segment.prototype.getField = function (idx, rep) {
    if (idx == 0) return this.name

    if (this.fields.length >= idx) {
        var f = this.fields[idx - 1]
        if (rep) {
            if (f.value.length >= rep) {
                return f.value[rep - 1].toString(delimiters)
            }

            return ""
        }

        return f.toString(delimiters)
    }

    return ""
}

segment.prototype.getComponent = function (idx, comp, sub) {
    if (this.fields.length >= idx) {
        var components = this.fields[idx - 1].value[0];
        if (components.length >= comp) {
            var target = components[comp - 1]
            if (sub) {
                if (target.value[0].length >= sub) {
                    return target.value[0][sub - 1].toString(delimiters)
                }

                return ""
            }

            return target.toString(delimiters)
        }

        return ""
    }

    return ""
}

segment.prototype.setComponent = function (idx, comp, val) {
    if (this.fields.length >= idx) {
        var components = this.fields[idx - 1].value[0];
        if (components.length >= comp) {
            components[idx - 1] = new component(val)
        }
    }
}

segment.prototype.toString = function (delimiters) {
    var res = this.name + delimiters.fieldSeparator
    for (var i = 0; i < this.fields.length; i++) {
        res += this.fields[i].toString(delimiters)
        if (i != this.fields.length - 1) res += delimiters.fieldSeparator
    }

    return res
}

function header() {
    this.name = "MSH";
    this.delimiters = {
        fieldSeparator: "|",
        componentSeparator: "^",
        subcomponentSeparator: "&",
        escapeCharacter: "\\",
        repetitionCharacter: "~",
        segmentSeparator: '\r'
    };

    this.fields = [];

    if (arguments.length > 1) {
        for (var i = 0; i < arguments.length; i++) {
            if (Array.isArray(arguments[i])) {
                var fs = new Array()
                for (var j = 0; j < arguments[i].length; j++) {
                    fs.push(new field(arguments[i][j]))
                }

                this.fields.push(fs)
            } else {
                this.fields.push(new field(arguments[i]))
            }
        }
    }
}

header.prototype.addField = segment.prototype.addField
header.prototype.setField = segment.prototype.setField
header.prototype.removeField = segment.prototype.removeField
header.prototype.getField = function (idx, rep) {
    if (idx == 0) return this.name
    if (idx == 1) return delimiters.fieldSeparator

    if (this.fields.length >= idx) {
        var f = this.fields[idx - 1]
        if (rep) {
            if (f.value.length >= rep) {
                return f.value[rep - 1].toString(delimiters)
            }

            return ""
        }

        return f.toString(delimiters)
    }

    return ""
}
header.prototype.getComponent = segment.prototype.getComponent

header.prototype.toString = function () {
    var res = this.name +
        this.delimiters.fieldSeparator +
        this.delimiters.componentSeparator +
        this.delimiters.repetitionCharacter +
        this.delimiters.escapeCharacter +
        this.delimiters.subcomponentSeparator +
        this.delimiters.fieldSeparator

    for (var i = 2; i < this.fields.length; i++) {
        res += this.fields[i].toString(this.delimiters)
        if (i != this.fields.length - 1) res += this.delimiters.fieldSeparator
    }

    return res
}


function message() {
    this.header = new header();
    this.segments = [];

    if (arguments.length > 0) {
        for (var i = 0; i < arguments.length; i++) {
            this.header.addField(arguments[i])
        }
    }
}

message.prototype.getSegment = function (name) {
    if (name == "MSH") return this.header

    for (var i = 0; i < this.segments.length; i++) {
        if (this.segments[i].name == name) return this.segments[i]
    }

    return null;
}

message.prototype.getSegments = function (name) {
    var res = [];
    for (var i = 0; i < this.segments.length; i++) {
        if (this.segments[i].name == name) res.push(this.segments[i])
    }

    return res
}

message.prototype.addSegment = function () {
    if (arguments.length == 1) {
        var s = new segment(arguments[0]);
        this.segments.push(s);
        return s;
    }

    if (arguments.length > 1) {
        var s = new segment(arguments[0]);
        for (var i = 1; i < arguments.length; i++) {
            s.addField(arguments[i]);
        }

        this.segments.push(s);
        return s;
    }
}

message.prototype.log = function () {
    var curSep = this.header.delimiters.segmentSeparator
    this.header.delimiters.segmentSeparator = '\n'
    var res = this.toString()

    this.header.delimiters.segmentSeparator = curSep

    return res
}

message.prototype.toString = function () {
    var res = this.header.toString() + this.header.delimiters.segmentSeparator;
    for (var i = 0; i < this.segments.length; i++) {
        res += this.segments[i].toString(this.header.delimiters);
        if (i != this.segments.length - 1) res += this.header.delimiters.segmentSeparator;
    }

    return res.replace(/^\s+|\s+$/g, '');
}

function parser(opts) {
    this.message = null;
    this.delimiters = {
        fieldSeparator: "|",
        componentSeparator: "^",
        subcomponentSeparator: "&",
        escapeCharacter: "\\",
        repetitionCharacter: "~",
        segmentSeparator: opts ? opts.segmentSeparator : '\r'
    }
}

parser.prototype.parse = function (s) {
    this.message = new message();
    var segments = s.split(this.delimiters.segmentSeparator);

    for (var i = 0; i < segments.length; i++) {
        if (i == 0) {
            this.message.header = this.parseHeader(segments[i]);
        } else {
            if (segments[i].trim() != "") {
                this.message.segments.push(this.parseSegment(segments[i]))
            }
        }
    }

    return this.message
}

parser.prototype.parseHeader = function (s) {
    var h = new header();
    var fields = s.split(this.delimiters.fieldSeparator);
    h.fields.push(new field(this.delimiters.fieldSeparator))

    for (var i = 0; i < fields.length; i++) {
        if (i == 0) {
            continue
        }
        if (i == 1) {
            h.fields.push(new field(fields[i]))
        } else {
            h.fields.push(this.parseField(fields[i]))
        }
    }

    return h
}

parser.prototype.parseSegment = function (s) {
    var seg = new segment()
    var fields = s.split(this.delimiters.fieldSeparator)

    seg.name = fields[0]

    for (var i = 1; i < fields.length; i++) {
        seg.fields.push(this.parseField(fields[i]))
    }

    return seg
}

parser.prototype.parseField = function (s) {
    var f = new field()
    if (!s) return f

    if (s.indexOf(this.delimiters.repetitionCharacter) != -1) {
        var reps = s.split(this.delimiters.repetitionCharacter)
        for (var i = 0; i < reps.length; i++) {
            f.value.push(this.parseField(reps[i]))
        }
    } else {
        var comps = s.split(this.delimiters.componentSeparator);
        var cs = []
        for (var i = 0; i < comps.length; i++) {
            cs.push(this.parseComponent(comps[i]))
        }

        f.value.push(cs)
    }

    return f;
}

parser.prototype.parseComponent = function (s) {
    var c = new component()
    if (s.indexOf(this.delimiters.repetitionCharacter) != -1) {
        var cs = s.split(this.delimiters.repetitionCharacter);
        for (var i = 0; i < cs.length; i++) {
            c.value.push(this.parseComponent(cs[i]));
        }
    } else {
        if (s.indexOf(this.delimiters.subcomponentSeparator) != -1) {
            var ss = s.split(this.delimiters.subcomponentSeparator);
            var subs = [];

            for (var i = 0; i < ss.length; i++) {
                subs.push(ss[i]);
            }

            c.value.push(subs);
        } else {
            c.value.push(s)
        }
    }

    return c;
}



// script body
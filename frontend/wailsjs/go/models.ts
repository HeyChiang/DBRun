export namespace connect {
	
	export class Config {
	    id: number;
	    type: string;
	    label: string;
	    username: string;
	    password: string;
	    host: string;
	    port: number;
	    database: string;
	    instance: string;
	    options: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.label = source["label"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.database = source["database"];
	        this.instance = source["instance"];
	        this.options = source["options"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace metadata {
	
	export class Credentials {
	    id: number;
	    type: string;
	    label: string;
	    username: string;
	    password: string;
	    host: string;
	    port: number;
	    database: string;
	    instance: string;
	    options: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Credentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.label = source["label"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.database = source["database"];
	        this.instance = source["instance"];
	        this.options = source["options"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class Style {
	    color: string;
	    isShow: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Style(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.color = source["color"];
	        this.isShow = source["isShow"];
	    }
	}
	export class Display {
	    dbCnt: number;
	    style: Record<string, Style>;
	
	    static createFrom(source: any = {}) {
	        return new Display(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dbCnt = source["dbCnt"];
	        this.style = this.convertValues(source["style"], Style, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ViewInfoVO {
	    id: number;
	    databaseId: number;
	    schemaId?: number;
	    alias: string;
	    name: string;
	    color: string;
	    definition: string;
	    remark: string;
	
	    static createFrom(source: any = {}) {
	        return new ViewInfoVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.databaseId = source["databaseId"];
	        this.schemaId = source["schemaId"];
	        this.alias = source["alias"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.definition = source["definition"];
	        this.remark = source["remark"];
	    }
	}
	export class FieldInfoVO {
	    id: number;
	    tableId: number;
	    alias: string;
	    name: string;
	    type: string;
	    display: boolean;
	    nullable?: boolean;
	    key: string;
	    comment?: string;
	    default_value?: string;
	    remark: string;
	    fontColor: string;
	    bgColor: string;
	    sort: number;
	
	    static createFrom(source: any = {}) {
	        return new FieldInfoVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tableId = source["tableId"];
	        this.alias = source["alias"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.display = source["display"];
	        this.nullable = source["nullable"];
	        this.key = source["key"];
	        this.comment = source["comment"];
	        this.default_value = source["default_value"];
	        this.remark = source["remark"];
	        this.fontColor = source["fontColor"];
	        this.bgColor = source["bgColor"];
	        this.sort = source["sort"];
	    }
	}
	export class TableInfoVO {
	    id: number;
	    databaseId: number;
	    schemaId?: number;
	    alias: string;
	    name: string;
	    color: string;
	    comment: string;
	    remark: string;
	    fields: FieldInfoVO[];
	
	    static createFrom(source: any = {}) {
	        return new TableInfoVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.databaseId = source["databaseId"];
	        this.schemaId = source["schemaId"];
	        this.alias = source["alias"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.comment = source["comment"];
	        this.remark = source["remark"];
	        this.fields = this.convertValues(source["fields"], FieldInfoVO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SchemaVO {
	    id: number;
	    databaseId: number;
	    alias: string;
	    name: string;
	    tables: TableInfoVO[];
	    views: ViewInfoVO[];
	
	    static createFrom(source: any = {}) {
	        return new SchemaVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.databaseId = source["databaseId"];
	        this.alias = source["alias"];
	        this.name = source["name"];
	        this.tables = this.convertValues(source["tables"], TableInfoVO);
	        this.views = this.convertValues(source["views"], ViewInfoVO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DatabaseInfoVO {
	    id: number;
	    configId: number;
	    alias: string;
	    name: string;
	    comment: string;
	    schemas: SchemaVO[];
	    tables: TableInfoVO[];
	    views: ViewInfoVO[];
	
	    static createFrom(source: any = {}) {
	        return new DatabaseInfoVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.configId = source["configId"];
	        this.alias = source["alias"];
	        this.name = source["name"];
	        this.comment = source["comment"];
	        this.schemas = this.convertValues(source["schemas"], SchemaVO);
	        this.tables = this.convertValues(source["tables"], TableInfoVO);
	        this.views = this.convertValues(source["views"], ViewInfoVO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DBInfoVO {
	    dbs: DatabaseInfoVO[];
	    display: Display;
	
	    static createFrom(source: any = {}) {
	        return new DBInfoVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dbs = this.convertValues(source["dbs"], DatabaseInfoVO);
	        this.display = this.convertValues(source["display"], Display);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	

}

export namespace service {
	
	export class TableCacheVO {
	    remark: string;
	    version: number;
	    fieldMap: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new TableCacheVO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.remark = source["remark"];
	        this.version = source["version"];
	        this.fieldMap = source["fieldMap"];
	    }
	}

}

export namespace sqlite {
	
	export class Project {
	    id: number;
	    name: string;
	    path: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


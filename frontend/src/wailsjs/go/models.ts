export namespace data {
	
	export class Task {
	    id: number;
	    name: string;
	    description: string;
	    done: boolean;
	    priority: number;
	    favorite: boolean;
	    inactivated: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.done = source["done"];
	        this.priority = source["priority"];
	        this.favorite = source["favorite"];
	        this.inactivated = source["inactivated"];
	    }
	}
	export class TimeEntry {
	    id: number;
	    taskID: number;
	    // Go type: time.Time
	    start: any;
	    // Go type: time.Time
	    end?: any;
	    synced: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TimeEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.taskID = source["taskID"];
	        this.start = this.convertValues(source["start"], null);
	        this.end = this.convertValues(source["end"], null);
	        this.synced = source["synced"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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


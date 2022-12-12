export namespace main {
	
	export class Task {
	    id: number;
	    name: string;
	    description: string;
	    done: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.done = source["done"];
	    }
	}

}


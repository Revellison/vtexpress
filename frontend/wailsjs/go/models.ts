export namespace models {
	
	export class VTSummary {
	    verdict: string;
	    totalEngines: number;
	    malicious: number;
	    suspicious: number;
	    undetected: number;
	    harmless: number;
	    timeout: number;
	    failure: number;
	    typeUnsupported: number;
	    fileName?: string;
	    fileType?: string;
	    sha256?: string;
	    size?: number;
	
	    static createFrom(source: any = {}) {
	        return new VTSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.verdict = source["verdict"];
	        this.totalEngines = source["totalEngines"];
	        this.malicious = source["malicious"];
	        this.suspicious = source["suspicious"];
	        this.undetected = source["undetected"];
	        this.harmless = source["harmless"];
	        this.timeout = source["timeout"];
	        this.failure = source["failure"];
	        this.typeUnsupported = source["typeUnsupported"];
	        this.fileName = source["fileName"];
	        this.fileType = source["fileType"];
	        this.sha256 = source["sha256"];
	        this.size = source["size"];
	    }
	}
	export class DomainInfo {
	    domain: string;
	    resolved?: string;
	    asn?: string;
	    provider?: string;
	    region?: string;
	
	    static createFrom(source: any = {}) {
	        return new DomainInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domain = source["domain"];
	        this.resolved = source["resolved"];
	        this.asn = source["asn"];
	        this.provider = source["provider"];
	        this.region = source["region"];
	    }
	}
	export class IPInfo {
	    address: string;
	    asn?: string;
	    provider?: string;
	    region?: string;
	
	    static createFrom(source: any = {}) {
	        return new IPInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.asn = source["asn"];
	        this.provider = source["provider"];
	        this.region = source["region"];
	    }
	}
	export class IPDomainReport {
	    ip: IPInfo[];
	    domain: DomainInfo[];
	
	    static createFrom(source: any = {}) {
	        return new IPDomainReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = this.convertValues(source["ip"], IPInfo);
	        this.domain = this.convertValues(source["domain"], DomainInfo);
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
	export class VTAnalysisEngine {
	    engineName: string;
	    category: string;
	    result?: string;
	    method?: string;
	    engineUpdate?: string;
	    engineVersion?: string;
	
	    static createFrom(source: any = {}) {
	        return new VTAnalysisEngine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engineName = source["engineName"];
	        this.category = source["category"];
	        this.result = source["result"];
	        this.method = source["method"];
	        this.engineUpdate = source["engineUpdate"];
	        this.engineVersion = source["engineVersion"];
	    }
	}
	export class EnginesVerdict {
	    engines: VTAnalysisEngine[];
	
	    static createFrom(source: any = {}) {
	        return new EnginesVerdict(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engines = this.convertValues(source["engines"], VTAnalysisEngine);
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
	export class ScanPayload {
	    engines_verdict: EnginesVerdict;
	    ip_domain: IPDomainReport;
	    vt_summary: VTSummary;
	
	    static createFrom(source: any = {}) {
	        return new ScanPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engines_verdict = this.convertValues(source["engines_verdict"], EnginesVerdict);
	        this.ip_domain = this.convertValues(source["ip_domain"], IPDomainReport);
	        this.vt_summary = this.convertValues(source["vt_summary"], VTSummary);
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
	export class AISummaryRequest {
	    payload: ScanPayload;
	
	    static createFrom(source: any = {}) {
	        return new AISummaryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.payload = this.convertValues(source["payload"], ScanPayload);
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
	export class AnalyzeByContentRequest {
	    fileName: string;
	    base64Data: string;
	    runAiSummary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeByContentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileName = source["fileName"];
	        this.base64Data = source["base64Data"];
	        this.runAiSummary = source["runAiSummary"];
	    }
	}
	export class AnalyzeByPathRequest {
	    path: string;
	    runAiSummary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeByPathRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.runAiSummary = source["runAiSummary"];
	    }
	}
	
	
	export class HistoryItem {
	    id: number;
	    createdAt: string;
	    fileName: string;
	    fileSha256: string;
	    payload: ScanPayload;
	    rawVt: number[];
	    aiSummary?: string;
	
	    static createFrom(source: any = {}) {
	        return new HistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = source["createdAt"];
	        this.fileName = source["fileName"];
	        this.fileSha256 = source["fileSha256"];
	        this.payload = this.convertValues(source["payload"], ScanPayload);
	        this.rawVt = source["rawVt"];
	        this.aiSummary = source["aiSummary"];
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
	
	
	export class SaveSettingsRequest {
	    vtApiKey: string;
	    aiApiKey: string;
	    aiProvider: string;
	    aiModel: string;
	    autoAiSummary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SaveSettingsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vtApiKey = source["vtApiKey"];
	        this.aiApiKey = source["aiApiKey"];
	        this.aiProvider = source["aiProvider"];
	        this.aiModel = source["aiModel"];
	        this.autoAiSummary = source["autoAiSummary"];
	    }
	}
	export class SaveSettingsResponse {
	    vtValid: boolean;
	    aiValid: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveSettingsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vtValid = source["vtValid"];
	        this.aiValid = source["aiValid"];
	        this.message = source["message"];
	    }
	}
	
	export class ScanResult {
	    fileName: string;
	    fileSha256: string;
	    scannedAt: string;
	    payload: ScanPayload;
	    rawVt: number[];
	    aiSummary?: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileName = source["fileName"];
	        this.fileSha256 = source["fileSha256"];
	        this.scannedAt = source["scannedAt"];
	        this.payload = this.convertValues(source["payload"], ScanPayload);
	        this.rawVt = source["rawVt"];
	        this.aiSummary = source["aiSummary"];
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
	export class SettingsState {
	    hasVtApiKey: boolean;
	    hasAiApiKey: boolean;
	    aiProvider: string;
	    aiModel: string;
	    autoAiSummary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SettingsState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasVtApiKey = source["hasVtApiKey"];
	        this.hasAiApiKey = source["hasAiApiKey"];
	        this.aiProvider = source["aiProvider"];
	        this.aiModel = source["aiModel"];
	        this.autoAiSummary = source["autoAiSummary"];
	    }
	}
	

}


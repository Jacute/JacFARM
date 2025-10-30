import { useEffect, useState } from "react"
import type { ExploitShort, LogLevel, Module } from "../../models/models"
import { getShortExploits } from "../../api/exploits"
import { Selector } from "../Selector"
import { getLogLevels, getModules } from "../../api/logs"

interface props {
    moduleId: number | null
    setModuleId: (module_id: number | null) => void
    exploitId: string | null
    setExploitId: (exploit_id: string | null) => void
    logLevelId: number | null
    setLogLevelId: (log_level_id: number | null) => void
    loadLogs: () => void
}

export const LogFilter = (props: props) => {
    const [modules, setModules] = useState<Module[]>([]);
    const [exploits, setExploits] = useState<ExploitShort[]>([]);
    const [logLevels, setLogLevels] = useState<LogLevel[]>([]);

    useEffect(() => {
        getShortExploits().then(exploits => setExploits(exploits));
        getModules().then(obj => setModules(obj.modules));
        getLogLevels().then(obj => setLogLevels(obj.log_levels));
    }, []);

    return (
        <div className="row mb-2">
          <Selector
            label="Module"
            value={props.moduleId}
            setValue={props.setModuleId}
            options={modules?.map((module) => ({ label: module.name, value: module.id }))}
          ></Selector>
          <Selector
            label="Exploit"
            value={props.exploitId}
            setValue={props.setExploitId}
            options={exploits?.map((exploit) => ({ label: exploit.name, value: exploit.id }))}
          ></Selector>
          <Selector
            label="Log Level"
            value={props.logLevelId}
            setValue={props.setLogLevelId}
            options={logLevels?.map((logLevel) => ({ label: logLevel.name, value: logLevel.id }))}
          ></Selector>
        <div className="col-3 d-flex align-items-end">
          <button className="btn btn-primary border border-1" onClick={props.loadLogs}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
              <path d="M23 4v6h-6"></path>
              <path d="M1 20v-6h6"></path>
              <path d="M2.51 9a9 9 0 0 1 14.85-3.36L23 10"></path>
              <path d="M21.49 15a9 9 0 0 1-14.85 3.36L1 14"></path>
            </svg>
          </button>
        </div>
      </div>
    )
}
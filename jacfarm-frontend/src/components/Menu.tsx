import { configPage } from "../pages/Config"
import { exploitsPage } from "../pages/Exploits"
import { flagsPage } from "../pages/Flags"
import { teamsPage } from "../pages/Teams"

interface props {
    setPage: (page: string) => void
}

export const Menu = (props: props) => {
    return (
        <div className="col-2 bg-primary border-end border-1 border-secondary m-0 p-0">
            <div className="logo h-25 d-flex align-items-center">
                <div className="row m-auto">
                    <div className="col-4 m-0 p-0">
                        <img src="/img/logo.jpg" className="rounded-2" style={{ width: "100px", height: "100px"}} />
                    </div>
                    <div className="col-8 d-flex align-items-center">
                        <h1 className="nice-header">
                            JacFARM
                        </h1>
                    </div>
                </div>
            </div>

            <div className="row d-flex mx-auto">
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(flagsPage)}>Flags</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(exploitsPage)}>Exploits</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(teamsPage)}>Teams</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(configPage)}>Config</button>
            </div>
        </div>
    )
}
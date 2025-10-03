import { Page } from '../pages/enum';
import type { PageType } from '../pages/enum';

interface props {
    flagsCount: number
    setPage: (page: PageType) => void
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
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(Page.FLAGS_PAGE)}>Flags</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(Page.EXPLOITS_PAGE)}>Exploits</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(Page.TEAMS_PAGE)}>Teams</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(Page.CONFIG_PAGE)}>Config</button>
                <button className="btn btn-primary w-90 border border-1 mb-2 mx-auto" onClick={() => props.setPage(Page.LOGS_PAGE)}>Logs</button>
            </div>

            <div className='row d-flex'>
                <p className='position-absolute bottom-0 start-0'>Flags in queue: {props.flagsCount}</p>
            </div>
        </div>
    )
}
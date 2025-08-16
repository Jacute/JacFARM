export const Menu = () => {
    return (
        <div className="col-2 bg-primary border-end border-5 border-secondary m-0 p-0">
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
                <button className="btn btn-primary w-100 border border-1 mb-2">Flags</button>
                <button className="btn btn-primary w-100 border border-1 mb-2">Exploits</button>
            </div>
        </div>
    )
}
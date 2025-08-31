import { useCallback, useState, useEffect } from "react";
import { ConfigTable } from "../components/config/Table";
import { Paginator } from "../components/Paginator";
import { getConfig } from "../api/config";
import type { Config } from "../models/models";
import { toast } from "react-toastify";

export const configPage = "config";

export const ConfigPage = () => {
    const [page, setPage] = useState<number>(1);
    const [count, setCount] = useState<number>(0);
    const [config, setConfig] = useState<Array<Config>>([]);

    const loadConfigs = useCallback(async () => {
        try {
            const [config, count] = await getConfig(page);
            setCount(count);
            setConfig(config);
        } catch (error: any) {
            toast.error(`Error loading configs: ${error.message}`);
            return;
        }
    }, [page]);

    useEffect(() => {
        loadConfigs();
    }, [page, loadConfigs]);

    return (
        <>
            <div className="row h-10 bg-primary border-bottom border-1 border-secondary"></div>

            <div className="row h-90">
                <div className="d-flex flex-column h-100">
                    <div className="row flex-grow-1 overflow-auto">
                        <ConfigTable
                            page={page}
                            configs={config}
                            setConfigs={setConfig}
                        />
                    </div>

                    <div className="row mt-auto mx-auto">
                        <div className="col">
                            <Paginator page={page} setPage={setPage} count={count} />
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};

import { toast, ToastContainer } from 'react-toastify';
import { Menu } from './components/Menu'
import { ExploitsPage } from './pages/Exploits';
import { FlagsPage } from './pages/Flags';
import './styles/index.css'
import { useEffect, useState } from 'react'
import { Page } from './pages/enum';
import type { PageType } from './pages/enum';
import { LogsPage } from './pages/Logs';
import { ConfigPage } from './pages/Config';
import { TeamsPage } from './pages/Teams';
import { getFlagsCount } from './api/flags';

function App() {
  const [page, setPage] = useState<PageType>(Page.FLAGS_PAGE);
  const [flagsCount, setFlagsCount] = useState<number>(0);

  const loadFlagsCount = async () => {
    try {
      const count = await getFlagsCount();
      setFlagsCount(count);
    } catch (error: any) {
      toast.error(`Error loading flags count: ${error.message}`);
    }
  };

  useEffect(() => {
    loadFlagsCount();

    const interval = setInterval(loadFlagsCount, 1000);

    return () => clearInterval(interval); // очистка при размонтировании
  }, [loadFlagsCount]);

  return (
    <div className="container-fluid vh-100">
      <div className="row h-100">
        <Menu setPage={setPage} flagsCount={flagsCount}/>
        <div className="col-10">
          {page === Page.FLAGS_PAGE && <FlagsPage />}
          {page === Page.EXPLOITS_PAGE && <ExploitsPage />}
          {page === Page.TEAMS_PAGE && <TeamsPage />}
          {page === Page.CONFIG_PAGE && <ConfigPage />}
          {page === Page.LOGS_PAGE && <LogsPage />}
        </div>
      </div>
      <ToastContainer limit={3} />
    </div>
  )
}

export default App;

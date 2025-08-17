import { Menu } from './components/Menu'
import { FlagsPage } from './pages/Flags';
import './styles/index.css'
import { useState } from 'react'

function App() {
  const [page, setPage] = useState<string>("flags");

  return (
    <div className="container-fluid vh-100">
      <div className="row h-100">
        <Menu setPage={setPage}/>
        <div className="col-10">
          <div className="row h-10 bg-primary border-bottom border-1 border-secondary">
            <div className="col-11"></div>
            <div className="col-1 d-flex align-items-center justify-content-center">
              <button className='btn btn-white w-50 border border-1'>
                ➕
              </button>
            </div>
          </div>
          <div className="row h-90">
            <FlagsPage />
          </div>
        </div>
      </div>
    </div>
  )
}

export default App

import { HashRouter, Routes, Route, Navigate } from "react-router-dom";

import Navbar from "src/components/Navbar";
import Footer from "src/components/Footer";
import Toaster from "src/components/Toast/Toaster";

import About from "src/screens/About";
import AppsList from "src/screens/apps/AppsList";
import ShowApp from "src/screens/apps/ShowApp";
import NewApp from "src/screens/apps/NewApp";
import AppCreated from "src/screens/apps/AppCreated";
import NotFound from "src/screens/NotFound";
import { useInfo } from "./hooks/useInfo";
import Loading from "./components/Loading";
import { Setup } from "./screens/Setup";

function App() {
  const { data: info } = useInfo();

  if (!info) {
    return <Loading />;
  }

  return (
    <div className="bg:white dark:bg-black min-h-full">
      <Toaster />
      <HashRouter>
        <Routes>
          <Route path="/" element={<Navbar />}>
            <Route
              index
              element={
                <Navigate
                  to={info.setupCompleted ? "/apps" : "/setup"}
                  replace
                />
              }
            />
            <Route path="setup" element={<Setup />} />
            <Route path="apps" element={<AppsList />} />
            <Route path="apps/:pubkey" element={<ShowApp />} />
            <Route path="apps/new" element={<NewApp />} />
            <Route path="apps/created" element={<AppCreated />} />
            <Route path="about" element={<About />} />
          </Route>
          <Route path="/*" element={<NotFound />} />
        </Routes>
        <Footer />
      </HashRouter>
    </div>
  );
}

export default App;

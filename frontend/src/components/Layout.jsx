import React from 'react';
import { Link } from 'react-router-dom';
import { Zap } from 'lucide-react';

const Layout = ({ children }) => {
  return (
    <div className="min-h-screen bg-[#0f172a] text-slate-100 font-sans selection:bg-indigo-500 selection:text-white overflow-hidden flex flex-col">
      {/* Dynamic Background Blurs */}
      <div className="fixed inset-0 z-0 pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-indigo-600/30 rounded-full mix-blend-screen filter blur-[100px] animate-blob"></div>
        <div className="absolute top-[20%] right-[-10%] w-[40%] h-[40%] bg-purple-600/30 rounded-full mix-blend-screen filter blur-[100px] animate-blob animation-delay-2000"></div>
        <div className="absolute bottom-[-20%] left-[20%] w-[50%] h-[50%] bg-pink-600/20 rounded-full mix-blend-screen filter blur-[120px] animate-blob animation-delay-4000"></div>
      </div>
      
      {/* Glassmorphic Header */}
      <header className="relative z-10 sticky top-0 bg-[#0f172a]/70 backdrop-blur-xl border-b border-white/10 shadow-lg">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-3 group transition-transform hover:scale-105">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500 p-[1px] shadow-lg shadow-indigo-500/20">
              <div className="w-full h-full bg-[#0f172a] rounded-xl flex items-center justify-center">
                <Zap className="w-5 h-5 text-indigo-400 group-hover:text-pink-400 transition-colors" />
              </div>
            </div>
            <h1 className="text-2xl font-black tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400">
              InciManage
            </h1>
          </Link>
          <div className="flex items-center gap-4">
             <div className="px-3 py-1.5 rounded-full bg-indigo-500/10 border border-indigo-500/20 flex items-center gap-2">
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                </span>
                <span className="text-xs font-medium text-emerald-400 tracking-wide uppercase">System Operational</span>
             </div>
          </div>
        </div>
      </header>

      <main className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full flex-grow flex flex-col">
        {children}
      </main>
    </div>
  );
};

export default Layout;

import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import { AlertCircle, Clock, ShieldAlert, Activity } from 'lucide-react';

const getSeverityStyles = (severity) => {
  switch (severity) {
    case 0: return 'bg-rose-500/10 text-rose-400 border-rose-500/20'; // P0 Critical
    case 1: return 'bg-orange-500/10 text-orange-400 border-orange-500/20'; // P1 High
    case 2: return 'bg-amber-500/10 text-amber-400 border-amber-500/20'; // P2 Medium
    default: return 'bg-blue-500/10 text-blue-400 border-blue-500/20'; // P3 Low
  }
};

const getSeverityLabel = (severity) => {
  switch (severity) {
    case 0: return 'SEV-0 CRITICAL';
    case 1: return 'SEV-1 HIGH';
    case 2: return 'SEV-2 MEDIUM';
    default: return 'SEV-3 LOW';
  }
};

const LiveFeed = () => {
  const [incidents, setIncidents] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchIncidents = async () => {
      try {
        const res = await axios.get('http://localhost:8080/api/v1/incidents');
        if (res.data) setIncidents(res.data);
      } catch (err) {
        console.error("Failed to fetch incidents", err);
      } finally {
        setLoading(false);
      }
    };
    
    fetchIncidents();
    const interval = setInterval(fetchIncidents, 5000); // Poll every 5s for live feel
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
         <div className="w-12 h-12 border-4 border-indigo-500/20 border-t-indigo-500 rounded-full animate-spin"></div>
         <p className="text-slate-400 font-medium animate-pulse">Scanning live signals...</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6 w-full animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex items-center justify-between mb-2">
        <h2 className="text-xl font-bold text-slate-100 flex items-center gap-2">
          <Activity className="text-indigo-400 w-5 h-5" />
          Active Incidents
        </h2>
        <span className="text-sm font-medium text-slate-400 bg-slate-800/50 px-3 py-1 rounded-full border border-slate-700">
          {incidents.length} Tracking
        </span>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {incidents.map((inc) => (
          <Link 
            to={`/incidents/${inc.id}`} 
            key={inc.id}
            className="group relative bg-[#1e293b]/60 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 hover:bg-[#1e293b]/80 hover:border-indigo-500/50 transition-all duration-300 hover:-translate-y-1 shadow-lg hover:shadow-indigo-500/10 flex flex-col gap-4 overflow-hidden"
          >
            {/* Glow effect on hover */}
            <div className="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
            
            <div className="flex justify-between items-start z-10 relative">
              <span className={`px-2.5 py-1 rounded-md border text-xs font-bold tracking-wider ${getSeverityStyles(inc.severity)}`}>
                {getSeverityLabel(inc.severity)}
              </span>
              <span className="px-2.5 py-1 rounded-md bg-slate-800 border border-slate-700 text-xs font-semibold text-slate-300 tracking-wider">
                {inc.status}
              </span>
            </div>

            <div className="z-10 relative">
              <h3 className="text-lg font-semibold text-slate-100 mb-1 group-hover:text-indigo-300 transition-colors">
                {inc.title}
              </h3>
              <p className="text-sm text-slate-400 line-clamp-2">
                {inc.description}
              </p>
            </div>

            <div className="mt-auto pt-4 flex items-center justify-between text-xs text-slate-500 border-t border-slate-700/50 z-10 relative">
              <div className="flex items-center gap-1.5">
                <ShieldAlert className="w-3.5 h-3.5" />
                <span className="truncate max-w-[100px]">{inc.component_id}</span>
              </div>
              <div className="flex items-center gap-1.5">
                <Clock className="w-3.5 h-3.5" />
                <span>{new Date(inc.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}</span>
              </div>
            </div>
          </Link>
        ))}
        {incidents.length === 0 && (
          <div className="col-span-full py-16 flex flex-col items-center justify-center text-slate-500 border border-dashed border-slate-700 rounded-2xl bg-slate-800/20">
            <ShieldAlert className="w-12 h-12 mb-3 text-slate-600" />
            <p className="font-medium">All clear. No active incidents detected.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default LiveFeed;

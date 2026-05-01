import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { ChevronLeft, Terminal, AlertTriangle, CheckCircle, Save, Calendar } from 'lucide-react';

const IncidentDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // RCA Form State
  const [rcaForm, setRcaForm] = useState({
    root_cause_category: 'SOFTWARE_BUG',
    fix_applied: '',
    prevention_steps: '',
    incident_start: '',
    incident_end: ''
  });
  const [submittingRCA, setSubmittingRCA] = useState(false);

  useEffect(() => {
    fetchDetail();
  }, [id]);

  const fetchDetail = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`http://localhost:8080/api/v1/incidents/${id}`);
      setData(res.data);
      // Pre-fill incident_start if available
      if (res.data.incident?.created_at) {
         setRcaForm(prev => ({...prev, incident_start: new Date(res.data.incident.created_at).toISOString().slice(0, 16)}));
      }
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  };

  const updateStatus = async (newStatus) => {
    try {
      await axios.patch(`http://localhost:8080/api/v1/incidents/${id}/status`, { status: newStatus });
      fetchDetail();
    } catch (err) {
      alert(`Transition Failed: ${err.response?.data?.error || err.message}`);
    }
  };

  const handleRCASubmit = async (e) => {
    e.preventDefault();
    setSubmittingRCA(true);
    try {
      // Convert local datetime-local string to RFC3339 for Go
      const payload = {
        ...rcaForm,
        incident_start: new Date(rcaForm.incident_start).toISOString(),
        incident_end: new Date(rcaForm.incident_end).toISOString(),
      };
      await axios.post(`http://localhost:8080/api/v1/incidents/${id}/rca`, payload);
      alert("RCA Submitted Successfully!");
      fetchDetail();
    } catch (err) {
      alert(`RCA Submission Failed: ${err.response?.data?.error || err.message}`);
    } finally {
      setSubmittingRCA(false);
    }
  };

  if (loading) return <div className="text-indigo-400 p-8 text-center animate-pulse">Decrypting signals...</div>;
  if (error) return <div className="text-rose-400 p-8 text-center bg-rose-500/10 rounded-xl border border-rose-500/20">{error}</div>;
  if (!data) return null;

  const { incident, raw_signals } = data;

  return (
    <div className="flex flex-col gap-6 animate-in fade-in zoom-in-95 duration-300">
      {/* Header */}
      <div className="flex items-center justify-between">
        <button onClick={() => navigate('/')} className="flex items-center gap-2 text-slate-400 hover:text-indigo-400 transition-colors">
          <ChevronLeft className="w-5 h-5" />
          <span className="font-medium">Back to Feed</span>
        </button>
        <div className="flex gap-2">
           {incident.status === 'OPEN' && (
             <button onClick={() => updateStatus('INVESTIGATING')} className="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg font-medium shadow-lg shadow-indigo-600/20 transition-all active:scale-95">Start Investigation</button>
           )}
           {incident.status === 'INVESTIGATING' && (
             <button onClick={() => updateStatus('RESOLVED')} className="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg font-medium shadow-lg shadow-emerald-600/20 transition-all active:scale-95">Mark Resolved</button>
           )}
           {incident.status === 'RESOLVED' && (
             <button onClick={() => updateStatus('CLOSED')} className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg font-medium shadow-lg transition-all active:scale-95">Close Ticket (Requires RCA)</button>
           )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Col: Details & RCA */}
        <div className="lg:col-span-2 flex flex-col gap-6">
          {/* Incident Info */}
          <div className="bg-[#1e293b]/80 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 shadow-xl">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h2 className="text-2xl font-bold text-slate-100">{incident.title}</h2>
                <p className="text-indigo-400 font-mono text-sm mt-1">{incident.component_id}</p>
              </div>
              <span className="px-3 py-1 rounded-md bg-slate-800 border border-slate-600 text-sm font-bold tracking-wider text-slate-200">
                {incident.status}
              </span>
            </div>
            <p className="text-slate-300 bg-slate-900/50 p-4 rounded-xl border border-slate-800">{incident.description}</p>
          </div>

          {/* RCA Form */}
          <div className="bg-[#1e293b]/80 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 shadow-xl relative overflow-hidden">
             <div className="absolute top-0 left-0 w-1 h-full bg-gradient-to-b from-pink-500 to-purple-500"></div>
             <h3 className="text-lg font-bold text-slate-100 mb-6 flex items-center gap-2">
                <AlertTriangle className="w-5 h-5 text-pink-400" />
                Root Cause Analysis
             </h3>
             <form onSubmit={handleRCASubmit} className="flex flex-col gap-5">
               <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                 <div className="flex flex-col gap-1.5">
                   <label className="text-sm font-medium text-slate-400 flex items-center gap-1.5"><Calendar className="w-4 h-4"/> Incident Start</label>
                   <input type="datetime-local" required value={rcaForm.incident_start} onChange={e => setRcaForm({...rcaForm, incident_start: e.target.value})} className="bg-slate-900/80 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500" />
                 </div>
                 <div className="flex flex-col gap-1.5">
                   <label className="text-sm font-medium text-slate-400 flex items-center gap-1.5"><Calendar className="w-4 h-4"/> Incident End</label>
                   <input type="datetime-local" required value={rcaForm.incident_end} onChange={e => setRcaForm({...rcaForm, incident_end: e.target.value})} className="bg-slate-900/80 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500" />
                 </div>
               </div>

               <div className="flex flex-col gap-1.5">
                 <label className="text-sm font-medium text-slate-400">Root Cause Category</label>
                 <select value={rcaForm.root_cause_category} onChange={e => setRcaForm({...rcaForm, root_cause_category: e.target.value})} className="bg-slate-900/80 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-indigo-500">
                    <option value="SOFTWARE_BUG">Software Bug</option>
                    <option value="HARDWARE_FAILURE">Hardware Failure</option>
                    <option value="NETWORK_ISSUE">Network Issue</option>
                    <option value="CONFIGURATION_ERROR">Configuration Error</option>
                    <option value="THIRD_PARTY">Third-Party Outage</option>
                 </select>
               </div>

               <div className="flex flex-col gap-1.5">
                 <label className="text-sm font-medium text-slate-400">Fix Applied</label>
                 <textarea required rows={3} value={rcaForm.fix_applied} onChange={e => setRcaForm({...rcaForm, fix_applied: e.target.value})} className="bg-slate-900/80 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-indigo-500 resize-none font-mono text-sm placeholder:text-slate-600" placeholder="Describe the immediate mitigation..."></textarea>
               </div>

               <div className="flex flex-col gap-1.5">
                 <label className="text-sm font-medium text-slate-400">Prevention Steps</label>
                 <textarea required rows={3} value={rcaForm.prevention_steps} onChange={e => setRcaForm({...rcaForm, prevention_steps: e.target.value})} className="bg-slate-900/80 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-indigo-500 resize-none font-mono text-sm placeholder:text-slate-600" placeholder="Long term fixes to prevent recurrence..."></textarea>
               </div>

               <button type="submit" disabled={submittingRCA} className="mt-2 w-full py-2.5 bg-gradient-to-r from-pink-600 to-purple-600 hover:from-pink-500 hover:to-purple-500 text-white font-bold rounded-lg shadow-lg shadow-pink-600/20 transition-all active:scale-[0.98] disabled:opacity-50 flex items-center justify-center gap-2">
                 <Save className="w-5 h-5" />
                 {submittingRCA ? 'Committing...' : 'Submit RCA to Database'}
               </button>
             </form>
          </div>
        </div>

        {/* Right Col: Raw Signals */}
        <div className="bg-[#0f172a]/90 backdrop-blur-md border border-slate-800 rounded-2xl flex flex-col overflow-hidden shadow-2xl h-[800px]">
          <div className="bg-slate-900 px-4 py-3 border-b border-slate-800 flex items-center justify-between">
            <h3 className="text-sm font-bold text-slate-300 flex items-center gap-2">
              <Terminal className="w-4 h-4 text-emerald-400" />
              Raw Signals (MongoDB)
            </h3>
            <span className="text-xs font-mono text-emerald-400 bg-emerald-400/10 px-2 py-0.5 rounded">{raw_signals?.length || 0} Events</span>
          </div>
          <div className="flex-1 overflow-auto p-4 font-mono text-xs text-slate-300 space-y-4">
            {raw_signals && raw_signals.map((sig, idx) => (
              <div key={idx} className="bg-slate-900/50 p-3 rounded-lg border border-slate-800">
                 <div className="flex justify-between text-slate-500 mb-2">
                   <span>{new Date(sig.timestamp).toISOString()}</span>
                   <span className="text-rose-400">{sig.error_code}</span>
                 </div>
                 <pre className="overflow-x-auto whitespace-pre-wrap text-emerald-300/80">
                   {JSON.stringify(sig.metadata, null, 2)}
                 </pre>
              </div>
            ))}
            {(!raw_signals || raw_signals.length === 0) && (
              <div className="text-slate-600 italic text-center py-10">No raw payloads found in data lake.</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default IncidentDetail;

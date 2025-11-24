import { useState, useEffect } from 'react'
import { Activity, HardDrive, FileText, CheckCircle, XCircle, Clock, RefreshCw } from 'lucide-react'

function App() {
  const [downloads, setDownloads] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [lastUpdated, setLastUpdated] = useState(new Date())

  const fetchDownloads = async () => {
    try {
      // Ensure your Go API server is running on port 8080
      const response = await fetch('http://localhost:8080/api/downloads')
      if (!response.ok) {
        throw new Error('Failed to fetch data')
      }
      const data = await response.json()
      // Handle case where API returns null for empty list
      setDownloads(data || [])
      setError(null)
    } catch (err) {
      console.error(err)
      setError('Could not connect to API Server. Is it running?')
    } finally {
      setLoading(false)
      setLastUpdated(new Date())
    }
  }

  // Initial fetch and set up interval
  useEffect(() => {
    fetchDownloads()
    const interval = setInterval(fetchDownloads, 2000) // Poll every 2 seconds
    return () => clearInterval(interval)
  }, [])

  // Helper to format date
  const formatDate = (dateString) => {
    if (!dateString) return 'N/A'
    return new Date(dateString).toLocaleString()
  }

  // Helper for status badges
  const getStatusColor = (status) => {
    switch (status) {
      case 'COMPLETED': return 'bg-green-100 text-green-800 border-green-200'
      case 'FAILED': return 'bg-red-100 text-red-800 border-red-200'
      case 'MISSING': return 'bg-yellow-100 text-yellow-800 border-yellow-200'
      default: return 'bg-gray-100 text-gray-800 border-gray-200'
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 font-sans p-8">
      <div className="max-w-6xl mx-auto">
        
        {/* Header Section */}
        <header className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
              <Activity className="text-blue-600" size={32} />
              KafkaSync
            </h1>
            <p className="text-gray-500 mt-1">Real-time file synchronization monitor</p>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-500 flex items-center gap-2 justify-end">
              <Clock size={14} />
              Last updated: {lastUpdated.toLocaleTimeString()}
            </div>
            <div className="mt-2 flex gap-4 text-sm font-medium text-gray-600">
              <div className="flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
                System Online
              </div>
              <div className="flex items-center gap-1">
                <HardDrive size={14} />
                {downloads.length} Files Processed
              </div>
            </div>
          </div>
        </header>

        {/* Main Content Area */}
        {error ? (
          <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
            <XCircle className="mx-auto text-red-500 mb-2" size={48} />
            <h3 className="text-lg font-medium text-red-800">Connection Error</h3>
            <p className="text-red-600">{error}</p>
            <button 
              onClick={fetchDownloads}
              className="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
            >
              Retry Connection
            </button>
          </div>
        ) : (
          <div className="bg-white shadow-sm border border-gray-200 rounded-xl overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">ID</th>
                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Filename</th>
                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Remote Path</th>
                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Status</th>
                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Timestamp</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {loading && downloads.length === 0 ? (
                    <tr>
                      <td colSpan="5" className="px-6 py-12 text-center text-gray-500">
                        <RefreshCw className="mx-auto animate-spin mb-2" />
                        Loading download history...
                      </td>
                    </tr>
                  ) : downloads.length === 0 ? (
                    <tr>
                      <td colSpan="5" className="px-6 py-12 text-center text-gray-500">
                        No downloads found in database.
                      </td>
                    </tr>
                  ) : (
                    downloads.map((item) => (
                      <tr key={item.id} className="hover:bg-gray-50 transition-colors">
                        <td className="px-6 py-4 text-sm text-gray-500">#{item.id}</td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            <FileText size={16} className="text-blue-500" />
                            <span className="font-medium text-gray-900">{item.filename}</span>
                          </div>
                          <div className="text-xs text-gray-400 mt-0.5 font-mono">Hash: {item.hash.substring(0, 8)}...</div>
                        </td>
                        <td className="px-6 py-4 text-sm text-gray-500 font-mono">{item.remote_location}</td>
                        <td className="px-6 py-4">
                          <span className={`px-3 py-1 rounded-full text-xs font-medium border ${getStatusColor(item.status)}`}>
                            {item.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 text-sm text-gray-500">
                          {formatDate(item.downloaded_at)}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
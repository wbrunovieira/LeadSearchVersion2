// src/App.tsx
import { useState } from 'react';
import LeadList from './components/LeadList';
import { categories } from './utils/categories';

interface Lead {
  [key: string]: any;
}

function App() {
  const [categoryID, setCategoryID] = useState<string>('');
  const [zipcodeID, setZipcodeID] = useState<string>('');
  const [radius, setRadius] = useState<number>(3000);
  const [maxResults, setMaxResults] = useState<number>(5);
  const [country, setCountry] = useState<string>('br'); // Estado para o país
  const [message, setMessage] = useState<string>('');
  const [messageType, setMessageType] = useState<'success' | 'error' | ''>('');
  const [leads, setLeads] = useState<Lead[]>([]);

  const BACKEND_URL_Search_GOOGLE = 'http://192.168.0.9:8082';
  const BACKEND_URL_API = 'http://192.168.0.9:8085';

  const handleStartSearch = async () => {
    try {
      setMessage('Iniciando busca...');
      setMessageType('');

      const normalizedZipcode = zipcodeID.replace(/\D/g, '');

      const url = `${BACKEND_URL_Search_GOOGLE}/start-search?category_id=${categoryID}&zipcode_id=${normalizedZipcode}&radius=${radius}&max_results=${maxResults}&country=${country}`;
      const response = await fetch(url);
      console.log('url', url);
      console.log('zipcodeID', zipcodeID);
      console.log('normalizedZipcode', normalizedZipcode);

      if (response.ok) {
        const text = await response.text();
        setMessage(`Busca concluída com sucesso! ${text}`);
        setMessageType('success');
        handleGetLeads();
      } else {
        setMessage('Erro ao iniciar a busca');
        setMessageType('error');
      }
    } catch (error) {
      console.error(error);
      setMessage('Erro de conexão com o backend');
      setMessageType('error');
    }
  };

  const handleGetLeads = async () => {
    try {
      const url = `${BACKEND_URL_API}/list-leads`;
      const response = await fetch(url);
      console.log('leads list response', response);
      if (response.ok) {
        const data: Lead[] = await response.json();
        console.log('leads list data', data);
        setLeads(data);
      } else {
        setMessage('Erro ao buscar leads');
        setMessageType('error');
      }
    } catch (error) {
      console.error(error);
      setMessage('Erro de conexão ao buscar leads');
      setMessageType('error');
    }
  };

  return (
    <div className="bg-slate-100 text-black p-4 max-w-4xl mx-auto my-8 rounded-lg shadow-md">
      <h1 className="text-3xl font-bold text-blue-600 mb-6">
        Buscar Leads no Google Places
      </h1>

      {/* Campo de seleção para categoria */}
      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Categoria:
        </label>
        <select
          value={categoryID}
          onChange={e => setCategoryID(e.target.value)}
          className="w-full p-2 border border-gray-300 rounded text-black"
        >
          <option value="">Selecione uma categoria</option>
          {categories.map(cat => (
            <option key={cat} value={cat}>
              {cat}
            </option>
          ))}
        </select>
      </div>

      {/* Campo para Zipcode */}
      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Zipcode ID:
        </label>
        <input
          type="text"
          value={zipcodeID}
          onChange={e => setZipcodeID(e.target.value)}
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      {/* Campo para seleção de País */}
      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          País:
        </label>
        <select
          value={country}
          onChange={e => setCountry(e.target.value)}
          className="w-full p-2 border border-gray-300 rounded text-black"
        >
          <option value="br">Brasil</option>
          <option value="us">Estados Unidos</option>
          <option value="ca">Canadá</option>
          {/* Adicione outras opções de país conforme necessário */}
        </select>
      </div>

      {/* Campo para Radius */}
      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Radius (m):
        </label>
        <input
          type="number"
          value={radius}
          onChange={e => setRadius(Number(e.target.value))}
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      {/* Campo para Max Results */}
      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Max Results:
        </label>
        <input
          type="number"
          value={maxResults}
          onChange={e =>
            setMaxResults(Number(e.target.value))
          }
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      <button
        onClick={handleStartSearch}
        className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-6 rounded-md font-semibold"
      >
        Iniciar Busca
      </button>

      {message && (
        <div
          className={`my-4 p-3 rounded-md ${
            messageType === 'success'
              ? 'bg-green-50 text-green-800 border border-green-300'
              : messageType === 'error'
              ? 'bg-red-50 text-red-800 border border-red-300'
              : 'bg-blue-50 text-blue-800 border border-blue-300'
          }`}
        >
          {messageType === 'success' && (
            <span className="inline-block mr-2">✅</span>
          )}
          {messageType === 'error' && (
            <span className="inline-block mr-2">❌</span>
          )}
          {message}
        </div>
      )}

      <hr className="my-8 border-gray-300" />

      <h2 className="text-2xl font-bold text-blue-600 mb-4">
        Leads Salvos
      </h2>

      <button
        onClick={handleGetLeads}
        className="bg-green-500 hover:bg-green-600 text-white py-2 px-6 rounded-md font-semibold mb-4"
      >
        Listar Leads
      </button>

      <LeadList leads={leads} />
    </div>
  );
}

export default App;
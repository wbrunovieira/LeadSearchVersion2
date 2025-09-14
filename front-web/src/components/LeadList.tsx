// /src/components/LeadList.tsx
import React from 'react';
import { LeadFront } from '../interfaces/leads';

interface LeadListProps {
  leads: LeadFront[];
}

const LeadList: React.FC<LeadListProps> = ({ leads }) => {
  if (!leads || leads.length === 0) {
    return <p className="text-gray-500">Nenhum lead encontrado.</p>;
  }

  const formatLabel = (key: string): string => {
    const labels: { [key: string]: string } = {
      businessName: 'Nome do Negócio',
      registeredName: 'Razão Social',
      cnpj: 'CNPJ',
      email: 'E-mail',
      phone: 'Telefone',
      whatsapp: 'WhatsApp',
      website: 'Website',
      instagram: 'Instagram',
      facebook: 'Facebook',
      tiktok: 'TikTok',
      address: 'Endereço',
      city: 'Cidade',
      state: 'Estado',
      zipCode: 'CEP',
      owner: 'Proprietário',
      category: 'Categoria',
      rating: 'Avaliação',
      description: 'Descrição',
      foundationDate: 'Data de Fundação'
    };
    return labels[key] || key;
  };

  const formatValue = (value: any): string => {
    if (value === null || value === undefined || value === '') {
      return 'N/A';
    }
    if (typeof value === 'object') {
      return JSON.stringify(value);
    }
    return String(value);
  };

  return (
    <div className="space-y-4">
      {leads.map((lead, index) => (
        <div
          key={lead.id || index}
          className="bg-white border border-gray-200 rounded-lg shadow-sm p-6"
        >
          <h3 className="text-xl font-bold text-blue-600 mb-4">
            {lead.businessName || `Lead #${index + 1}`}
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {/* Informações Principais */}
            <div className="space-y-2">
              <h4 className="font-semibold text-gray-700 border-b pb-1">Informações Principais</h4>
              <p><span className="font-medium">CNPJ:</span> {formatValue(lead.cnpj)}</p>
              <p><span className="font-medium">Razão Social:</span> {formatValue(lead.registeredName)}</p>
              <p><span className="font-medium">Proprietário:</span> {formatValue(lead.owner)}</p>
              <p><span className="font-medium">Categoria:</span> {formatValue(lead.category)}</p>
              <p><span className="font-medium">Avaliação:</span> {lead.rating > 0 ? `${lead.rating} ⭐` : 'N/A'}</p>
              <p><span className="font-medium">Data de Fundação:</span> {formatValue(lead.foundationDate)}</p>
            </div>

            {/* Contato */}
            <div className="space-y-2">
              <h4 className="font-semibold text-gray-700 border-b pb-1">Contato</h4>
              <p><span className="font-medium">Telefone:</span> {formatValue(lead.phone)}</p>
              <p><span className="font-medium">WhatsApp:</span> {formatValue(lead.whatsapp)}</p>
              <p><span className="font-medium">E-mail:</span> {formatValue(lead.email)}</p>
              <p><span className="font-medium">Website:</span> {lead.website ? (
                <a href={lead.website} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                  {lead.website}
                </a>
              ) : 'N/A'}</p>
            </div>

            {/* Localização */}
            <div className="space-y-2">
              <h4 className="font-semibold text-gray-700 border-b pb-1">Localização</h4>
              <p><span className="font-medium">Endereço:</span> {formatValue(lead.address)}</p>
              <p><span className="font-medium">Cidade:</span> {formatValue(lead.city)}</p>
              <p><span className="font-medium">Estado:</span> {formatValue(lead.state)}</p>
              <p><span className="font-medium">CEP:</span> {formatValue(lead.zipCode)}</p>
            </div>

            {/* Redes Sociais */}
            <div className="space-y-2">
              <h4 className="font-semibold text-gray-700 border-b pb-1">Redes Sociais</h4>
              <p><span className="font-medium">Instagram:</span> {lead.instagram ? (
                <a href={lead.instagram} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                  {lead.instagram}
                </a>
              ) : 'N/A'}</p>
              <p><span className="font-medium">Facebook:</span> {lead.facebook ? (
                <a href={lead.facebook} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                  {lead.facebook}
                </a>
              ) : 'N/A'}</p>
              <p><span className="font-medium">TikTok:</span> {formatValue(lead.tiktok)}</p>
            </div>
          </div>

          {/* Descrição */}
          {lead.description && lead.description !== '' && (
            <div className="mt-4 pt-4 border-t">
              <h4 className="font-semibold text-gray-700 mb-2">Descrição</h4>
              <p className="text-gray-600">{lead.description}</p>
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default LeadList;

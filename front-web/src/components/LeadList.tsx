// /src/components/LeadList.tsx
import React from 'react';

interface Lead {
  [key: string]: any;
}

interface LeadListProps {
  leads: Lead[];
}

const LeadList: React.FC<LeadListProps> = ({ leads }) => {
  if (!leads || leads.length === 0) {
    return <p>Nenhum lead encontrado.</p>;
  }

  return (
    <div>
      {leads.map((lead, index) => (
        <div
          key={index}
          style={{
            border: '1px solid #ccc',
            borderRadius: '4px',
            padding: '1rem',
            marginBottom: '1rem',
          }}
        >
          <h3>Lead #{index + 1}</h3>
          <ul style={{ listStyle: 'none', padding: 0 }}>
            {Object.entries(lead).map(([key, value]) => (
              <li
                key={key}
                style={{ marginBottom: '0.5rem' }}
              >
                <strong>{key}:</strong>{' '}
                {typeof value === 'object'
                  ? JSON.stringify(value)
                  : value !== null
                  ? value.toString()
                  : 'null'}
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
};

export default LeadList;

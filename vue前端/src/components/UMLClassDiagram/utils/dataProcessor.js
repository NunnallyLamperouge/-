export const processNodes = (nodes) => {
  return nodes.map((node, index) => ({
    id: node.id,
    name: node.name,
    symbolSize: 50,
    category: 0,
    label: {
      show: true,
      formatter: (params) => {
        const node = nodes.find(n => n.id === params.data.id);
        return [
          `{name|${node.name}}`,
          `{attributes|${node.attributes.join('\n')}}`,
          `{methods|${node.methods.join('\n')}}`
        ].join('\n');
      },
      rich: {
        name: {
          fontSize: 14,
          fontWeight: 'bold',
          padding: [5, 10, 5, 10],
          backgroundColor: '#f0f0f0',
          borderRadius: 4
        },
        attributes: {
          fontSize: 12,
          padding: [5, 10, 5, 10],
          color: '#666666'
        },
        methods: {
          fontSize: 12,
          padding: [5, 10, 5, 10],
          color: '#666666'
        }
      }
    },
    x: Math.cos(index * 2 * Math.PI / nodes.length) * 300,
    y: Math.sin(index * 2 * Math.PI / nodes.length) * 300
  }));
};

export const processEdges = (edges) => {
  return edges.map(edge => ({
    source: edge.source,
    target: edge.target,
    label: {
      show: true,
      formatter: edge.label,
      position: 'middle',
      backgroundColor: '#ffffff',
      padding: [4, 8],
      borderRadius: 4,
      fontSize: 12,
      color: '#666666'
    }
  }));
}; 
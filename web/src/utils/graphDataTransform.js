export function transformGraphData(backendData) {
  if (!backendData) {
    return { nodes: [], edges: [] }
  }

  const nodes = Array.isArray(backendData.nodes) ? backendData.nodes : []
  const edges = Array.isArray(backendData.edges) ? backendData.edges : []
  
  const nodeIdMap = new Map()

  const transformedNodes = nodes.map((node, index) => {
    const nodeId = node.entity_name || `node_${index}`
    const nodeLabel = node.entity_name || `node_${index}`
    const nodeSize = node.pagerank ? Math.max(40, Math.min(60, node.pagerank * 100)) : 50

    if (node && node.entity_name) {
      nodeIdMap.set(node.entity_name, nodeId)
    }
    if (node && node.entity_id) {
      nodeIdMap.set(String(node.entity_id), nodeId)
    }
    if (node && node.id) {
      nodeIdMap.set(String(node.id), nodeId)
    }

    // 排除原始 node 中的 style，让 graphConfig 的默认样式生效
    const { style, ...nodeWithoutStyle } = node || {}
    
    return {
      ...nodeWithoutStyle,
      id: nodeId,
      label: nodeLabel,
      originalLabel: nodeLabel,
      type: 'circle',
      size: nodeSize
    }
  })

  const transformedEdges = edges.map((edge, index) => {
    const edgeId = `e${index}`
    const edgeLabel = edge.description || ''

    const source =
      nodeIdMap.get(edge && edge.source_entity) ||
      nodeIdMap.get(edge && edge.source) ||
      nodeIdMap.get(
        edge && edge.source_id ? String(edge.source_id) : undefined
      ) ||
      (edge && edge.source_entity) ||
      (edge && edge.source) ||
      (edge && edge.source_id ? String(edge.source_id) : undefined) ||
      `source_${index}`

    const target =
      nodeIdMap.get(edge && edge.target_entity) ||
      nodeIdMap.get(edge && edge.target) ||
      nodeIdMap.get(
        edge && edge.target_id ? String(edge.target_id) : undefined
      ) ||
      (edge && edge.target_entity) ||
      (edge && edge.target) ||
      (edge && edge.target_id ? String(edge.target_id) : undefined) ||
      `target_${index}`

    // 排除原始 edge 中的 style，让 graphConfig 的默认样式生效
    const { style, ...edgeWithoutStyle } = edge || {}
    
    return {
      ...edgeWithoutStyle,
      id: edgeId,
      source,
      target,
      label: edgeLabel
    }
  })
  return {
    nodes: transformedNodes,
    edges: transformedEdges
  }
}

export default {
  transformGraphData
}

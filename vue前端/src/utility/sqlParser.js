/**
 * SQL解析工具
 * 将SQL语句解析为ER图所需的JSON对象
 */

/**
 * 解析SQL语句
 * @param {string} sqlText - SQL语句文本
 * @returns {Object} 解析后的表结构和关系
 */
export function parseSQL(sqlText) {
  if (!sqlText) return { tables: [], relationships: [] };

  // 移除SQL注释
  sqlText = removeComments(sqlText);
  
  // 提取CREATE TABLE语句
  const createTableStatements = extractCreateTableStatements(sqlText);
  
  // 解析表结构
  const tables = createTableStatements.map(parseCreateTable);
  
  // 提取表之间的关系
  const relationships = extractRelationships(tables);
  
  return {
    tables,
    relationships
  };
}

/**
 * 移除SQL注释
 */
function removeComments(sqlText) {
  // 移除单行注释 --
  sqlText = sqlText.replace(/--.*$/gm, '');
  
  // 移除多行注释 /* */
  sqlText = sqlText.replace(/\/\*[\s\S]*?\*\//g, '');
  
  return sqlText;
}

/**
 * 提取CREATE TABLE语句
 */
function extractCreateTableStatements(sqlText) {
  const regex = /CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?[`"']?(\w+)[`"']?\s*\(([\s\S]*?)(?:\)\s*(?:ENGINE|DEFAULT|CHARACTER|COLLATE|AUTO_INCREMENT|;))/gi;
  const statements = [];
  let match;
  
  while ((match = regex.exec(sqlText)) !== null) {
    statements.push({
      tableName: match[1],
      content: match[2]
    });
  }
  
  return statements;
}

/**
 * 解析CREATE TABLE语句
 */
function parseCreateTable(statement) {
  const { tableName, content } = statement;
  const lines = content.split(',').map(line => line.trim());
  const fields = [];
  const constraints = [];
  
  // 提取字段和约束
  lines.forEach(line => {
    if (line.startsWith('PRIMARY KEY') || line.startsWith('FOREIGN KEY') || 
        line.startsWith('CONSTRAINT') || line.startsWith('UNIQUE') || 
        line.startsWith('INDEX') || line.startsWith('KEY')) {
      constraints.push(line);
    } else {
      // 解析字段定义
      const fieldMatch = line.match(/^[`"']?(\w+)[`"']?\s+(\w+(?:\(\d+(?:,\d+)?\))?)/);
      if (fieldMatch) {
        const fieldName = fieldMatch[1];
        const fieldType = fieldMatch[2];
        const isPrimary = line.toUpperCase().includes('PRIMARY KEY');
        const isNotNull = line.toUpperCase().includes('NOT NULL');
        const isAutoIncrement = line.toUpperCase().includes('AUTO_INCREMENT');
        
        fields.push({
          name: fieldName,
          type: fieldType,
          isPrimary,
          isNotNull,
          isAutoIncrement,
          isForeign: false // 稍后处理外键
        });
      }
    }
  });
  
  // 处理主键约束
  const primaryKeyConstraint = constraints.find(c => c.startsWith('PRIMARY KEY'));
  if (primaryKeyConstraint) {
    const pkFields = extractFieldsFromConstraint(primaryKeyConstraint);
    pkFields.forEach(pkField => {
      const field = fields.find(f => f.name === pkField);
      if (field) {
        field.isPrimary = true;
      }
    });
  }
  
  // 处理外键约束
  const foreignKeyConstraints = constraints.filter(c => c.startsWith('FOREIGN KEY') || 
    (c.startsWith('CONSTRAINT') && c.includes('FOREIGN KEY')));
  
  const foreignKeys = foreignKeyConstraints.map(constraint => {
    const fkMatch = constraint.match(/FOREIGN KEY\s*\(\s*[`"']?(\w+)[`"']?\s*\)\s*REFERENCES\s*[`"']?(\w+)[`"']?\s*\(\s*[`"']?(\w+)[`"']?\s*\)/i);
    if (fkMatch) {
      return {
        field: fkMatch[1],
        referencedTable: fkMatch[2],
        referencedField: fkMatch[3]
      };
    }
    return null;
  }).filter(Boolean);
  
  // 标记外键字段
  foreignKeys.forEach(fk => {
    const field = fields.find(f => f.name === fk.field);
    if (field) {
      field.isForeign = true;
      field.references = {
        table: fk.referencedTable,
        field: fk.referencedField
      };
    }
  });
  
  return {
    name: tableName,
    type: '实体表',
    fields,
    foreignKeys
  };
}

/**
 * 从约束中提取字段名
 */
function extractFieldsFromConstraint(constraint) {
  const fieldsMatch = constraint.match(/\(\s*([^)]+)\s*\)/);
  if (fieldsMatch) {
    return fieldsMatch[1].split(',').map(field => {
      // 移除引号和空格
      return field.replace(/[`"'\s]/g, '');
    });
  }
  return [];
}

/**
 * 提取表之间的关系
 */
function extractRelationships(tables) {
  const relationships = [];
  
  tables.forEach(table => {
    if (table.foreignKeys && table.foreignKeys.length > 0) {
      table.foreignKeys.forEach(fk => {
        relationships.push({
          source: table.name,
          target: fk.referencedTable,
          sourceField: fk.field,
          targetField: fk.referencedField,
          type: 'one-to-many' // 默认为一对多关系
        });
      });
    }
  });
  
  return relationships;
}

/**
 * 分析表之间的关系类型
 * 根据外键约束推断关系类型（一对一、一对多、多对多）
 */
export function analyzeRelationshipTypes(parsedData) {
  const { tables, relationships } = parsedData;
  
  // 查找多对多关系（通过中间表）
  const tableMap = new Map(tables.map(table => [table.name, table]));
  
  // 复制关系数组，以便我们可以修改它
  const enhancedRelationships = [...relationships];
  
  // 查找可能的中间表（具有两个外键的表）
  tables.forEach(table => {
    if (table.foreignKeys && table.foreignKeys.length >= 2) {
      // 检查该表是否可能是中间表
      const nonKeyFields = table.fields.filter(field => 
        !field.isPrimary && !field.isForeign
      );
      
      // 如果表只有外键和主键，没有其他字段，可能是中间表
      if (nonKeyFields.length <= 1) {
        // 找到这个表相关的所有关系
        const tableRelationships = enhancedRelationships.filter(
          rel => rel.source === table.name
        );
        
        // 如果有两个以上的外键关系，可能是多对多
        if (tableRelationships.length >= 2) {
          for (let i = 0; i < tableRelationships.length; i++) {
            for (let j = i + 1; j < tableRelationships.length; j++) {
              const rel1 = tableRelationships[i];
              const rel2 = tableRelationships[j];
              
              // 添加一个新的多对多关系
              enhancedRelationships.push({
                source: rel1.target,
                target: rel2.target,
                sourceField: rel1.targetField,
                targetField: rel2.targetField,
                type: 'many-to-many',
                throughTable: table.name
              });
            }
          }
        }
      }
    }
  });
  
  // 检查一对一关系
  enhancedRelationships.forEach(rel => {
    if (rel.type === 'one-to-many') {
      const sourceTable = tableMap.get(rel.source);
      const targetTable = tableMap.get(rel.target);
      
      if (sourceTable && targetTable) {
        const sourceField = sourceTable.fields.find(f => f.name === rel.sourceField);
        
        // 如果外键也是主键，可能是一对一关系
        if (sourceField && sourceField.isPrimary) {
          rel.type = 'one-to-one';
        }
      }
    }
  });
  
  return {
    tables,
    relationships: enhancedRelationships
  };
} 
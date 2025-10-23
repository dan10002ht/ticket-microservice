import { db } from '../config/databaseConfig.js';
import IRepository from './interfaces/IRepository.js';

/**
 * Base Repository với Master-Slave Pattern
 * Tự động route operations dựa trên loại operation
 * Implements IRepository interface
 */
class BaseRepository extends IRepository {
  constructor(tableName) {
    super();
    this.tableName = tableName;
    // PgPool-II sẽ tự động route queries dựa trên loại operation:
    // - SELECT queries → slave databases
    // - INSERT/UPDATE/DELETE queries → master database
    // - Transactions → master database
    this.db = db(tableName);
  }

  // ========== READ OPERATIONS (Sử dụng Slave) ==========

  /**
   * Tìm tất cả records
   */
  async findAll(options = {}) {
    const { limit, offset, orderBy, orderDirection = 'asc' } = options;
    let query = this.db.select('*');

    if (orderBy) {
      query = query.orderBy(orderBy, orderDirection);
    }

    if (limit) {
      query = query.limit(limit);
    }

    if (offset) {
      query = query.offset(offset);
    }

    return await query;
  }

  /**
   * Tìm record theo ID (internal)
   */
  async findById(id) {
    const results = await this.db.select('*').where('id', id).limit(1);

    return results[0] || null;
  }

  /**
   * Tìm record theo điều kiện
   */
  async findOne(conditions) {
    const results = await this.db.select('*').where(conditions).limit(1);

    return results[0] || null;
  }

  /**
   * Tìm nhiều records theo điều kiện
   */
  async findMany(conditions, options = {}) {
    const { limit, offset, orderBy, orderDirection = 'asc' } = options;
    let query = this.db.select('*').where(conditions);

    if (orderBy) {
      query = query.orderBy(orderBy, orderDirection);
    }

    if (limit) {
      query = query.limit(limit);
    }

    if (offset) {
      query = query.offset(offset);
    }

    return await query;
  }

  /**
   * Đếm số records theo điều kiện
   */
  async count(conditions = {}) {
    const result = await this.db.count('* as total').where(conditions);

    return parseInt(result[0].total);
  }

  /**
   * Kiểm tra record có tồn tại không
   */
  async exists(conditions) {
    const result = await this.db.select('id').where(conditions).limit(1);

    return result.length > 0;
  }

  /**
   * Tìm record theo public_id (external)
   */
  async findByPublicId(publicId) {
    const results = await this.db.select('*').where('public_id', publicId).limit(1);

    return results[0] || null;
  }

  // ========== WRITE OPERATIONS (Sử dụng Master) ==========

  /**
   * Tạo record mới
   */
  async create(data) {
    const [result] = await this.db.insert(data).returning('*');

    return result;
  }

  /**
   * Tạo nhiều records cùng lúc
   */
  async createMany(dataArray) {
    const results = await this.db.insert(dataArray).returning('*');

    return results;
  }

  /**
   * Cập nhật record theo ID
   */
  async updateById(id, data) {
    const [result] = await this.db.where('id', id).update(data).returning('*');

    return result;
  }

  /**
   * Cập nhật records theo điều kiện
   */
  async update(conditions, data) {
    const results = await this.db.where(conditions).update(data).returning('*');

    return results;
  }

  /**
   * Cập nhật record theo public_id
   */
  async updateByPublicId(publicId, data) {
    const [result] = await this.db.where('public_id', publicId).update(data).returning('*');

    return result;
  }

  /**
   * Xóa record theo ID
   */
  async deleteById(id) {
    const [result] = await this.db.where('id', id).del().returning('*');

    return result;
  }

  /**
   * Xóa records theo điều kiện
   */
  async delete(conditions) {
    const results = await this.db.where(conditions).del().returning('*');

    return results;
  }

  /**
   * Xóa record theo public_id
   */
  async deleteByPublicId(publicId) {
    const [result] = await this.db.where('public_id', publicId).del().returning('*');

    return result;
  }

  /**
   * Upsert - Insert hoặc Update
   */
  async upsert(data, uniqueColumns = ['id']) {
    const conflictColumns = uniqueColumns.join(', ');
    const [result] = await this.db.insert(data).onConflict(conflictColumns).merge().returning('*');

    return result;
  }

  // ========== TRANSACTION OPERATIONS ==========

  /**
   * Thực hiện transaction với master database
   */
  async transaction(callback) {
    return await this.db.transaction(callback);
  }

  // ========== ADVANCED QUERY OPERATIONS ==========

  /**
   * Raw query - PgPool-II sẽ tự động route queries
   */
  async rawQuery(sql, params = []) {
    return await db.raw(sql, params);
  }

  /**
   * Join query với slave database
   */
  async join(joinTable, joinCondition, options = {}) {
    const { select = '*', where = {}, limit, offset, orderBy, orderDirection = 'asc' } = options;

    let query = this.db.select(select).join(joinTable, joinCondition).where(where);

    if (orderBy) {
      query = query.orderBy(orderBy, orderDirection);
    }

    if (limit) {
      query = query.limit(limit);
    }

    if (offset) {
      query = query.offset(offset);
    }

    return await query;
  }

  /**
   * Pagination helper
   */
  async paginate(page = 1, pageSize = 10, conditions = {}, options = {}) {
    const offset = (page - 1) * pageSize;
    const { orderBy, orderDirection = 'asc' } = options;

    const [data, totalCount] = await Promise.all([
      this.findMany(conditions, { limit: pageSize, offset, orderBy, orderDirection }),
      this.count(conditions),
    ]);

    return {
      data,
      pagination: {
        page,
        pageSize,
        total: totalCount,
        totalPages: Math.ceil(totalCount / pageSize),
        hasNext: page < Math.ceil(totalCount / pageSize),
        hasPrev: page > 1,
      },
    };
  }
}

export default BaseRepository;

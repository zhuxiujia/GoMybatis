<template>
  <div class="body">
    <MonacoEditor
      height="1000"
      language="xml"
      :code="code"
      :editorOptions="options"
    >
    </MonacoEditor>
  </div>

</template>

<script>
  import MonacoEditor from 'vue-monaco-editor';
  export default {
    data() {
      return {
        code: '<?xml version="1.0" encoding="UTF-8"?>\n' +
          '<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"\n' +
          '        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">\n' +
          '<mapper>\n' +
          '    <resultMap id="BaseResultMap">\n' +
          '        <id column="id" property="id"/>\n' +
          '        <result column="name" property="name" langType="string"/>\n' +
          '        <result column="pc_link" property="pcLink" langType="string"/>\n' +
          '        <result column="h5_link" property="h5Link" langType="string"/>\n' +
          '        <result column="remark" property="remark" langType="string"/>\n' +
          '        <result column="create_time" property="createTime" langType="time.Time"/>\n' +
          '        <result column="delete_flag" property="deleteFlag" langType="int"/>\n' +
          '    </resultMap>\n' +
          '    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->\n' +
          '    <!-- 后台查询产品 -->\n' +
          '    <select id="selectByCondition" resultMap="BaseResultMap">\n' +
          '        <bind name="pattern" value="\'%\' + name + \'%\'"/>\n' +
          '        select * from biz_activity\n' +
          '        <where>\n' +
          '            <if test="name != \'\'">\n' +
          '                <!--可以使用bind标签 and name like #{pattern}-->\n' +
          '                and name like #{pattern}\n' +
          '                <!--可以使用默认 and name like concat(\'%\',#{name},\'%\')-->\n' +
          '                <!--and name like concat(\'%\',#{name},\'%\')-->\n' +
          '            </if>\n' +
          '            <if test="startTime != \'\'">and create_time >= #{startTime}</if>\n' +
          '            <if test="endTime != \'\'">and create_time &lt;= #{endTime}</if>\n' +
          '        </where>\n' +
          '        order by create_time desc\n' +
          '        <if test="page >= 0 and size != 0">limit #{page}, #{size}</if>\n' +
          '    </select>\n' +
          '    <!--int countByCondition(@Param("name")String name,@Param("startTime") Date startTime, @Param("endTime")Date endTime);-->\n' +
          '    <select id="countByCondition">\n' +
          '        select count(id) from biz_activity where delete_flag=1\n' +
          '        <if test="name != \'\'">\n' +
          '            and name like concat(\'%\',#{name},\'%\')\n' +
          '        </if>\n' +
          '        <if test="startTime != \'\'">\n' +
          '            and create_time >= #{startTime}\n' +
          '        </if>\n' +
          '        <if test="endTime != \'\'">\n' +
          '            and create_time &lt;= #{endTime}\n' +
          '        </if>\n' +
          '    </select>\n' +
          '    <!--List<Activity> selectAll();-->\n' +
          '    <select id="selectAll">\n' +
          '        select * from biz_activity where delete_flag=1 order by create_time desc\n' +
          '    </select>\n' +
          '    <!--Activity selectByUUID(@Param("uuid")String uuid);-->\n' +
          '    <select id="selectByUUID">\n' +
          '        select * from biz_activity\n' +
          '        where uuid = #{uuid}\n' +
          '        and delete_flag = 1\n' +
          '    </select>\n' +
          '    <select id="selectById">\n' +
          '        select * from biz_activity\n' +
          '        where id = #{id}\n' +
          '        and delete_flag = 1\n' +
          '    </select>\n' +
          '    <select id="selectByIds">\n' +
          '        select * from biz_activity\n' +
          '        where delete_flag = 1\n' +
          '        <foreach separator="," collection="ids" item="item" index="index" open=" and id in (" close=")">\n' +
          '            #{item}\n' +
          '        </foreach>\n' +
          '    </select>\n' +
          '    <update id="deleteById">\n' +
          '        update biz_activity\n' +
          '        set delete_flag = 0\n' +
          '        where id = #{id}\n' +
          '    </update>\n' +
          '    <update id="updateById">\n' +
          '        update biz_activity\n' +
          '        <set>\n' +
          '            <if test="name != \'\'">name = #{name},</if>\n' +
          '            <if test="pcLink != \'\'">pc_link = #{pcLink},</if>\n' +
          '            <if test="h5Link != \'\'">h5_link = #{h5Link},</if>\n' +
          '            <if test="remark != \'\'">remark = #{remark},</if>\n' +
          '            <if test="createTime != \'\'">create_time = #{createTime},</if>\n' +
          '            <if test="deleteFlag != \'\'">delete_flag = #{deleteFlag},</if>\n' +
          '        </set>\n' +
          '        where id = #{id} and delete_flag = 1\n' +
          '    </update>\n' +
          '    <insert id="insert">\n' +
          '        insert into biz_activity\n' +
          '        <trim prefix="(" suffix=")" suffixOverrides=",">\n' +
          '            <if test="id != \'\'">id,</if>\n' +
          '            <if test="name != \'\'">name,</if>\n' +
          '            <if test="pcLink != \'\'">pc_link,</if>\n' +
          '            <if test="h5Link != \'\'">h5_link,</if>\n' +
          '            <if test="remark != \'\'">remark,</if>\n' +
          '            <if test="createTime != \'\'">create_time,</if>\n' +
          '            <if test="deleteFlag != \'\'">delete_flag,</if>\n' +
          '        </trim>\n' +
          '\n' +
          '        <trim prefix="values (" suffix=")" suffixOverrides=",">\n' +
          '            <if test="id != \'\'">#{id},</if>\n' +
          '            <if test="name != \'\'">#{name},</if>\n' +
          '            <if test="pcLink != \'\'">#{pcLink},</if>\n' +
          '            <if test="h5Link != \'\'">#{h5Link},</if>\n' +
          '            <if test="remark != \'\'">#{remark},</if>\n' +
          '            <if test="createTime != \'\'">#{createTime},</if>\n' +
          '            <if test="deleteFlag != \'\'">#{deleteFlag},</if>\n' +
          '        </trim>\n' +
          '    </insert>\n' +
          '\n' +
          '    <select id="choose" resultMap="BaseResultMap">\n' +
          '        SELECT * FROM biz_activity\n' +
          '        <choose>\n' +
          '            <when test="deleteFlag > 1">WHERE delete_flag > 1</when>\n' +
          '            <when test="deleteFlag == 1">WHERE delete_flag = 1</when>\n' +
          '            <otherwise>WHERE delete_flag = #{deleteFlag}</otherwise>\n' +
          '        </choose>\n' +
          '    </select>\n' +
          '\n' +
          '    <sql id="links">\n' +
          '        pc_link,h5_link\n' +
          '        <!--不启用TypeConvert的话，使用${} 而不是 #{}-->\n' +
          '        <if test="column != \'\'">,${column}</if>\n' +
          '    </sql>\n' +
          '\n' +
          '    <select id="selectLinks">\n' +
          '        select\n' +
          '        <include refid="links"/>\n' +
          '        from biz_activity where delete_flag = 1\n' +
          '    </select>\n' +
          '</mapper>',
        options: {
          selectOnLineNumbers: false
        }
      };
    },
    methods: {
      onMounted(editor) {
        this.editor = editor;
      },
      onCodeChange(editor) {
        console.log(editor.getValue());
      },
      openFiles(){

      },
    },
    mounted:function(){
      this.openFiles();
    },
    components: {
      MonacoEditor
    }
  }
</script>

<style lang="less">

  .body {
    height: auto;
  }


</style>
